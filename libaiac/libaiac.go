package libaiac

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/briandowns/spinner"
	"github.com/google/uuid"
	"github.com/ido50/requests"
	"github.com/manifoldco/promptui"
)

const (
	ChatGPTHost           = "chat.openai.com"
	DefaultUserAgent      = "Mozilla/5.0 (Windows NT 10.0; rv:107.0) Gecko/20100101 Firefox/107.0"
	SessionTokenCookie    = "__Secure-next-auth.session-token"
	CloudflareTokenCookie = "cf_clearance"
	CloudflareBmCookie    = "__cf_bm"
	CallbackURLCookie     = "__Secure-next-auth.callback-url"
)

var ErrNoCode = errors.New("no code generated")

type Client struct {
	*requests.HTTPClient
	token               string
	cloudflareClearance string
	cloudflareBm        string
	userAgent           string
	chatGPT             bool
}

type AIACClientInput struct {
	ChatGPT             bool
	Token               string
	CloudflareClearance string
	CloudflareBm        string
	UserAgent           string
}

func NewClient(input *AIACClientInput) *Client {
	cli := &Client{
		token:               input.Token,
		chatGPT:             input.ChatGPT,
		cloudflareClearance: input.CloudflareClearance,
		cloudflareBm:        input.CloudflareBm,
		userAgent:           input.UserAgent,
	}

	if !cli.chatGPT {
		cli.HTTPClient = requests.NewClient("https://api.openai.com/v1").
			Accept("application/json").
			Header("Authorization", fmt.Sprintf("Bearer %s", cli.token)).
			ErrorHandler(func(
				httpStatus int,
				contentType string,
				body io.Reader,
			) error {
				var res struct {
					Error struct {
						Message string `json:"message"`
						Type    string `json:"type"`
					} `json:"error"`
				}
				err := json.NewDecoder(body).Decode(&res)
				if err != nil {
					return fmt.Errorf(
						"OpenAI returned response %s",
						http.StatusText(httpStatus),
					)
				}

				return fmt.Errorf("[%s] %s", res.Error.Type, res.Error.Message)
			})
	} else {
		ua := cli.userAgent
		if len(ua) == 0 {
			ua = DefaultUserAgent
		}
		cli.HTTPClient = requests.NewClient(fmt.Sprintf("https://%s", ChatGPTHost)).
			Header("User-Agent", ua).
			Header("Accept-Language", "en-US,en;q=0.9")
	}

	return cli
}

func (client *Client) Ask(
	ctx context.Context,
	prompt string,
	shouldRetry bool,
	shouldQuit bool,
	outputPath string,
	readmePath string,
) (err error) {
	spin := spinner.New(spinner.CharSets[2],
		100*time.Millisecond,
		spinner.WithWriter(color.Error),
		spinner.WithSuffix("\tGenerating code ..."))
	spin.Start()
	killed := false

	defer func() {
		if !killed {
			spin.Stop()
		}
	}()

	var code, readme string

	if client.chatGPT {
		code, readme, err = client.askViaChatGPT(ctx, prompt)
	} else {
		code, err = client.askViaAPI(ctx, prompt)
	}

	if err != nil {
		return err
	}

	code = fmt.Sprintf("%s\n", code)

	spin.Stop()
	killed = true

	fmt.Fprintf(os.Stdout, code)
	if shouldQuit {
		return nil
	}
	if shouldRetry {
		input := promptui.Prompt{
			Label: "Hit [S/s] to save the file or [R/r] to retry [Q/q] to quit",
			Validate: func(s string) error {
				if strings.ToLower(s) != "s" && strings.ToLower(s) != "r" && strings.ToLower(s) != "q" {
					return fmt.Errorf("Invalid input. Try again please.")
				}
				return nil
			},
		}

		result, err := input.Run()

		if strings.ToLower(result) == "q" {
			// finish without saving
			return nil
		} else if err != nil || strings.ToLower(result) == "r" {
			// retry once more
			return client.Ask(ctx, prompt, shouldRetry, shouldQuit, outputPath, readmePath)
		}
	}

	if outputPath == "-" {
		input := promptui.Prompt{
			Label: "Enter a file path",
		}

		outputPath, err = input.Run()
		if err != nil {
			return err
		}
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf(
			"failed creating output file %s: %w",
			outputPath, err,
		)
	}

	defer f.Close()

	fmt.Fprint(f, code)

	if readmePath != "" {
		f, err := os.Create(readmePath)
		if err != nil {
			return fmt.Errorf(
				"failed creating readme file %s: %w",
				readmePath, err,
			)
		}

		defer f.Close()

		fmt.Fprintf(f, readme)
	}

	if outputPath != "-" {
		fmt.Printf("Code saved successfully at %s\n", outputPath)
	}

	return nil
}

func (client *Client) askViaChatGPT(ctx context.Context, prompt string) (
	code string,
	readme string,
	err error,
) {
	requestID, err := uuid.NewRandom()
	if err != nil {
		return code, readme, fmt.Errorf("failed generating UUID: %w", err)
	}

	accessToken, err := client.loadAccessToken(ctx)
	if err != nil {
		return code, readme, fmt.Errorf("failed loading access token: %w", err)
	}

	cacheAccessToken(accessToken)

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

	parentID, err := uuid.NewRandom()
	if err != nil {
		return code, readme, fmt.Errorf(
			"failed generating initial parent ID: %w",
			err,
		)
	}

	body["parent_message_id"] = parentID.String()

	messages := make(chan []byte)

	err = client.NewRequest("POST", "/backend-api/conversation").
		Accept("text/event-stream").
		JSONBody(body).
		Header("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Cookie(&http.Cookie{
			Name:   CallbackURLCookie,
			Value:  "https://" + ChatGPTHost + "/",
			Path:   "/",
			Domain: ChatGPTHost,
		}).
		Cookie(&http.Cookie{
			Name:   SessionTokenCookie,
			Value:  client.token,
			Path:   "/",
			Domain: ChatGPTHost,
		}).
		Cookie(&http.Cookie{
			Name:   CloudflareTokenCookie,
			Value:  client.cloudflareClearance,
			Path:   "/",
			Domain: ChatGPTHost,
		}).
		Cookie(&http.Cookie{
			Name:   CloudflareBmCookie,
			Value:  client.cloudflareBm,
			Path:   "/",
			Domain: ChatGPTHost,
		}).
		Subscribe(ctx, messages)
	if err != nil {
		return code, readme, fmt.Errorf("failed starting a conversation: %w", err)
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
			return code, readme, fmt.Errorf(
				"failed parsing ChatGPT response: %w",
				err,
			)
		}

		if len(msg.Message.Content.Parts) < 1 {
			continue
		}

		finalMessage = msg.Message.Content.Parts[0]
	}

	scanner := bufio.NewScanner(strings.NewReader(finalMessage))

	var writeOutput, alreadyHadCode bool
	var b strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if line == "```" {
			if !alreadyHadCode {
				writeOutput = !writeOutput
				if !writeOutput {
					alreadyHadCode = true
				}
			}
		} else if writeOutput {
			fmt.Fprintln(&b, line)
		}
	}

	return b.String(), finalMessage, nil
}

func (client *Client) askViaAPI(ctx context.Context, prompt string) (
	code string,
	err error,
) {
	var answer struct {
		Choices []struct {
			Text         string `json:"text"`
			Index        int64  `json:"index"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}

	var status int
	err = client.NewRequest("POST", "/completions").
		JSONBody(map[string]interface{}{
			"model":      "text-davinci-003",
			"prompt":     prompt,
			"max_tokens": 4097 - len(prompt),
		}).
		Into(&answer).
		StatusInto(&status).
		RunContext(ctx)
	if err != nil {
		return code, fmt.Errorf("failed sending prompt: %w", err)
	}

	if len(answer.Choices) == 0 {
		return code, fmt.Errorf("no results returned from API")
	}

	if answer.Choices[0].FinishReason != "stop" {
		return code, fmt.Errorf(
			"result was truncated by API due to %s",
			answer.Choices[0].FinishReason,
		)
	}

	return strings.TrimSpace(answer.Choices[0].Text), nil
}

type CacheFile struct {
	AccessToken string `json:"accessToken"`
	Expiry      int64  `json:"expiry"`
}

func (client *Client) loadAccessToken(ctx context.Context) (token string, err error) {
	// try to load access token from cache file
	f, err := os.Open(filepath.Join(xdg.ConfigHome, "aiac.token"))
	if err == nil {
		defer f.Close()

		var tokenData CacheFile
		err = json.NewDecoder(f).Decode(&tokenData)
		if err == nil && time.Now().Unix() < tokenData.Expiry {
			// cached token has not expired yet
			return tokenData.AccessToken, nil
		}
	}

	// get an access token from ChatGPT
	var session struct {
		AccessToken string `json:"accessToken"`
	}

	err = client.NewRequest("GET", "/api/auth/session").
		Cookie(&http.Cookie{
			Name:   SessionTokenCookie,
			Value:  client.token,
			Path:   "/",
			Domain: ChatGPTHost,
		}).
		Cookie(&http.Cookie{
			Name:   CloudflareTokenCookie,
			Value:  client.cloudflareClearance,
			Path:   "/",
			Domain: ChatGPTHost,
		}).
		Cookie(&http.Cookie{
			Name:   CloudflareBmCookie,
			Value:  client.cloudflareBm,
			Path:   "/",
			Domain: ChatGPTHost,
		}).
		Into(&session).
		RunContext(ctx)
	if err != nil {
		return token, fmt.Errorf(
			"failed getting session details: %w",
			err,
		)
	}

	return session.AccessToken, nil
}

func cacheAccessToken(token string) error {
	f, err := os.Create(filepath.Join(xdg.ConfigHome, "aiac.token"))
	if err != nil {
		return fmt.Errorf("failed creating token file: %w", err)
	}

	defer f.Close()

	err = json.NewEncoder(f).Encode(CacheFile{
		AccessToken: token,
		Expiry:      time.Now().Add(3600 * time.Second).Unix(),
	})
	if err != nil {
		return fmt.Errorf("failed encoding token data: %w", err)
	}

	return nil
}
