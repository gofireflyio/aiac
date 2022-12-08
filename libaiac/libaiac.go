package libaiac

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ido50/requests"
)

const (
	ChatGPTHost        = "chat.openai.com"
	DefaultUserAgent   = "Mozilla/5.0 (Windows NT 10.0; rv:107.0) Gecko/20100101 Firefox/107.0"
	SessionTokenCookie = "__Secure-next-auth.session-token"
	CallbackURLCookie  = "__Secure-next-auth.callback-url"
)

var ErrNoCode = errors.New("no code generated")

type Client struct {
	*requests.HTTPClient
	sessionToken   string
	accessToken    string
	conversationID string
	parentID       string
}

func NewClient(sessionToken string) *Client {
	sessionToken = strings.TrimPrefix(sessionToken, SessionTokenCookie+"=")

	return &Client{
		sessionToken: sessionToken,
		HTTPClient: requests.NewClient("https://"+ChatGPTHost).
			Accept("application/json").
			Header("User-Agent", DefaultUserAgent).
			Header("X-Openai-Assistant-App-Id", "").
			Header("Accept-Language", "en-US,en;q=0.9").
			Header("Referer", "https://"+ChatGPTHost+"/chat"),
	}
}

func (client *Client) Ask(
	ctx context.Context,
	prompt string,
	outputPath string,
	readmePath string,
) (err error) {
	requestID, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed generating UUID: %w", err)
	}

	messages := make(chan []byte)

	if client.accessToken == "" {
		// get an access token
		var session struct {
			AccessToken string `json:"accessToken"`
		}

		err = client.NewRequest("GET", "/api/auth/session").
			Cookie(&http.Cookie{
				Name:   SessionTokenCookie,
				Value:  client.sessionToken,
				Path:   "/",
				Domain: ChatGPTHost,
			}).
			Into(&session).
			RunContext(ctx)
		if err != nil {
			return fmt.Errorf("failed getting session details: %w", err)
		}

		client.accessToken = session.AccessToken
	}

	// start a conversation
	body := map[string]interface{}{
		"action": "next",
		"messages": []map[string]interface{}{
			{
				"id":   requestID.String(),
				"role": "user",
				"content": map[string]interface{}{
					"content_type": "text",
					"parts":        []string{prompt},
				},
			},
		},
		"model":           "text-davinci-002-render",
		"conversation_id": nil,
	}

	if client.conversationID != "" {
		body["conversation_id"] = client.conversationID
	}
	if client.parentID == "" {
		parentID, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("failed generating initial parent ID: %w", err)
		}

		client.parentID = parentID.String()
	}

	body["parent_message_id"] = client.parentID

	err = client.NewRequest("POST", "/backend-api/conversation").
		Accept("text/event-stream").
		JSONBody(body).
		Header("Authorization", fmt.Sprintf("Bearer %s", client.accessToken)).
		Cookie(&http.Cookie{
			Name:   CallbackURLCookie,
			Value:  "https://" + ChatGPTHost + "/",
			Path:   "/",
			Domain: ChatGPTHost,
		}).
		Cookie(&http.Cookie{
			Name:   SessionTokenCookie,
			Value:  client.sessionToken,
			Path:   "/",
			Domain: ChatGPTHost,
		}).
		Subscribe(ctx, messages)
	if err != nil {
		return fmt.Errorf("failed starting a conversation: %w", err)
	}

	type ChatGPTMessage struct {
		Message struct {
			ID      string `json:"id"`
			Content struct {
				ContentType string   `json:"content_type"`
				Parts       []string `json:"parts"`
			} `json:"content"`
		} `json:"message"`
		ConversationID string  `json:"conversation_id"`
		Error          *string `json:"error"`
	}

	var finalMessage string

	for bmsg := range messages {
		bmsg = bytes.TrimSpace(bytes.TrimPrefix(bmsg, []byte("data: ")))

		if bytes.Compare(bmsg, []byte{'[', 'D', 'O', 'N', 'E', ']'}) == 0 {
			break
		}

		var msg ChatGPTMessage
		err = json.Unmarshal(bmsg, &msg)
		if err != nil {
			return fmt.Errorf("failed parsing ChatGPT response: %w", err)
		}

		client.conversationID = msg.ConversationID
		client.parentID = msg.Message.ID

		if len(msg.Message.Content.Parts) < 1 {
			continue
		}

		finalMessage = msg.Message.Content.Parts[0]
	}

	var output io.Writer

	if outputPath == "-" {
		output = os.Stdout
	} else {
		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf(
				"failed creating output file %s: %w",
				outputPath, err,
			)
		}

		defer f.Close()

		output = f
	}

	var readme io.Writer

	if readmePath != "" {
		f, err := os.OpenFile(readmePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf(
				"failed creating/opening readme file %s: %w",
				readmePath, err,
			)
		}

		defer f.Close()

		readme = f
	}

	if readme != nil {
		fmt.Fprintf(readme, "# %s\n", prompt)
	}

	scanner := bufio.NewScanner(strings.NewReader(finalMessage))

	var writeOutput, alreadyHadCode bool

	for scanner.Scan() {
		line := scanner.Text()

		if readme != nil {
			fmt.Fprintln(readme, line)
		}

		if line == "```" {
			if !alreadyHadCode {
				writeOutput = !writeOutput
				if !writeOutput {
					alreadyHadCode = true
				}
			}
		} else if writeOutput {
			fmt.Fprintln(output, line)
		}
	}

	if readme != nil {
		fmt.Fprintln(readme, "")
	}

	if !alreadyHadCode {
		return ErrNoCode
	}

	return nil
}
