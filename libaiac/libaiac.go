package libaiac

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/ido50/requests"
)

// Version contains aiac's version string
var Version = "development"

// Client is a structure used to continuously generate IaC code via OpenAPI/ChatGPT
type Client struct {
	*requests.HTTPClient
	apiKey string
}

var (
	// ErrResultTruncated is returned when the OpenAI API returned a truncated
	// result. The reason for the truncation will be appended to the error
	// string.
	ErrResultTruncated = errors.New("result was truncated")

	// ErrNoResults is returned if the OpenAI API returned an empty result. This
	// should not generally happen.
	ErrNoResults = errors.New("no results return from API")

	// ErrUnsupportedModel is returned if the SetModel method is provided with
	// an unsupported model
	ErrUnsupportedModel = errors.New("unsupported model")

	// ErrUnexpectedStatus is returned when the OpenAI API returned a response
	// with an unexpected status code
	ErrUnexpectedStatus = errors.New("OpenAI returned unexpected response")

	// ErrRequestFailed is returned when the OpenAI API returned an error for
	// the request
	ErrRequestFailed = errors.New("request failed")
)

// NewClient creates a new instance of the Client struct, with the provided
// input options. Neither the OpenAI API nor ChatGPT are yet contacted at this
// point.
func NewClient(apiKey string) *Client {
	if apiKey == "" {
		return nil
	}

	cli := &Client{
		apiKey: strings.TrimPrefix(apiKey, "Bearer "),
	}

	cli.HTTPClient = requests.NewClient("https://api.openai.com/v1").
		Accept("application/json").
		Header("Authorization", fmt.Sprintf("Bearer %s", cli.apiKey)).
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
					"%w %s",
					ErrUnexpectedStatus,
					http.StatusText(httpStatus),
				)
			}

			return fmt.Errorf(
				"%w: [%s]: %s",
				ErrRequestFailed,
				res.Error.Type,
				res.Error.Message,
			)
		})

	return cli
}

// Response is the struct returned from methods generating code via the OpenAI
// API.
type Response struct {
	// FullOutput is the complete output returned by the API. This is generally
	// a Markdown-formatted message that contains the generated code, plus
	// explanations, if any.
	FullOutput string

	// Code is the extracted code section from the complete output. If code was
	// not found or extraction otherwise failed, this will be the same as
	// FullOutput.
	Code string

	// APIKeyUsed is the API key used when making the request.
	APIKeyUsed string

	// TokensUsed is the number of tokens utilized by the request. This is
	// the "usage.total_tokens" value returned from the API.
	TokensUsed int64
}

// GenerateCode sends the provided prompt to the OpenAI API and returns a
// Response object. It is a convenience wrapper around client.Complete (for
// text completion models) and client.Chat.Send (for chat models).
func (client *Client) GenerateCode(
	ctx context.Context,
	model Model,
	prompt string,
) (res Response, err error) {
	if model.Type == ModelTypeChat {
		chat := client.Chat(model)
		return chat.Send(ctx, prompt)
	}

	return client.Complete(ctx, model, prompt)
}

var codeRegex = regexp.MustCompile("(?ms)^```(?:[^\n]*)\n(.*?)\n```$")

// ExtractCode receives the full output string from the OpenAI API and attempts
// to extract a code block from it. OpenAI code blocks are generally Markdown
// blocks surrounded by the ``` string on both sides. If successful, the code
// string will be returned together with a true value, otherwise an empty string
// is returned together with a false value.
func ExtractCode(output string) (string, bool) {
	m := codeRegex.FindStringSubmatch(output)
	if m == nil || m[1] == "" {
		return "", false
	}

	return m[1], true
}
