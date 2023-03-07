package libaiac

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/ido50/requests"
	"github.com/manifoldco/promptui"
)

// Client is a structure used to continuously generate IaC code via OpenAPI/ChatGPT
type Client struct {
	*requests.HTTPClient
	apiKey string
	model  Model
	full   bool
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

// Model is an enum used to select the language model to use
type Model string

const (
	// ModelChatGPT represents the gpt-3.5-turbo model used by ChatGPT.
	ModelChatGPT = "gpt-3.5-turbo"

	// ModelTextDaVinci3 represents the text-davinci-003 language generation
	// model.
	ModelTextDaVinci3 = "text-davinci-003"

	// ModelCodeDaVinci2 represents the code-davinci-002 code generation model.
	ModelCodeDaVinci2 = "code-davinci-002"
)

// Decode is used by the kong library to map CLI-provided values to the Model
// type
func (m *Model) Decode(ctx *kong.DecodeContext) error {
	var provided string

	err := ctx.Scan.PopValueInto("string", &provided)
	if err != nil {
		return fmt.Errorf("failed getting model value: %w", err)
	}

	for _, supported := range []Model{
		ModelChatGPT,
		ModelTextDaVinci3,
		ModelCodeDaVinci2,
	} {
		if string(supported) == provided {
			*m = supported
			return nil
		}
	}

	return fmt.Errorf("%w %s", ErrUnsupportedModel, provided)
}

// SupportedModels is a list of all models supported by aiac
var SupportedModels = []string{ModelChatGPT, ModelTextDaVinci3, ModelCodeDaVinci2}

// MaxTokens is the maximum amount of tokens supported by the model used. Newer
// OpenAI models support a maximum of 4096 tokens.
var MaxTokens = 4096

// NewClient creates a new instance of the Client struct, with the provided
// input options. Neither the OpenAI API nor ChatGPT are yet contacted at this
// point.
func NewClient(apiKey string) *Client {
	if apiKey == "" {
		return nil
	}

	cli := &Client{
		apiKey: strings.TrimPrefix(apiKey, "Bearer "),
		model:  ModelChatGPT,
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

// SetModel changes the language model to use with the OpenAI API
func (client *Client) SetModel(model Model) *Client {
	client.model = model
	return client
}

// SetFull sets whether output is returned/stored in full, including
// explanations (if any), or if only the code is extracted. Defaults to false
// (meaning only code is extracted).
func (client *Client) SetFull(full bool) *Client {
	client.full = full
	return client
}

// Ask asks the OpenAI API to generate code based on the provided prompt.
// It is only meant to be used in command line applications (see GenerateCode
// for library usage). The generated code will always be printed to standard
// output, but may optionally be stored in the file whose path is provided by
// the outputPath argument. To only print to standard output, provide an empty
// string or a dash ("-") via outputPath. If shouldRetry is true, you will be
// prompted whether to regenerate the response after it is printed to standard output,
// in case you are unhappy with the response. If shouldQuit is true, the code
// is printed to standard output and the function returns, without storing to a
// file or asking whether to regenerate the response.
func (client *Client) Ask(
	ctx context.Context,
	prompt string,
	shouldRetry bool,
	shouldQuit bool,
	outputPath string,
) (err error) {
	spin := spinner.New(spinner.CharSets[2],
		100*time.Millisecond, //nolint: gomnd
		spinner.WithWriter(color.Error),
		spinner.WithSuffix("\tGenerating code ..."))

	spin.Start()

	killed := false

	defer func() {
		if !killed {
			spin.Stop()
		}
	}()

	code, err := client.GenerateCode(ctx, prompt)
	if err != nil {
		return err
	}

	spin.Stop()

	killed = true

	fmt.Fprintln(os.Stdout, code)

	if shouldQuit {
		return nil
	}

	if shouldRetry {
		errInvalidInput := errors.New("invalid input, please try again") //nolint: goerr113

		input := promptui.Prompt{
			Label: "Hit [S/s] to save the file, [R/r] to retry, [M/m] to modify prompt, [Q/q] to quit",
			Validate: func(s string) error {
				switch strings.ToLower(s) {
				case "s", "r", "m", "q":
					return nil
				}

				return errInvalidInput
			},
		}

		result, err := input.Run()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return fmt.Errorf("prompt failed: %w", err)
		}

		switch strings.ToLower(result) {
		case "r":
			// retry once more
			return client.Ask(ctx, prompt, shouldRetry, shouldQuit, outputPath)
		case "m":
			// let user modify prompt
			input := promptui.Prompt{
				Label:   "New prompt",
				Default: prompt,
			}

			prompt, err = input.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}

			return client.Ask(ctx, prompt, shouldRetry, shouldQuit, outputPath)
		case "q":
			// finish without saving
			return nil
		}
	}

	if outputPath == "" {
		input := promptui.Prompt{
			Label: "Enter a file path",
		}

		outputPath, err = input.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
	}

	if outputPath != "-" {
		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf(
				"failed creating output file %s: %w",
				outputPath, err,
			)
		}

		defer f.Close()

		fmt.Fprintln(f, code)

		fmt.Fprintf(os.Stderr, "Code saved successfully at %s\n", outputPath)
	}

	return nil
}

var codeRegex = regexp.MustCompile("(?ms)^```(?:[^\n]*)\n(.*?)\n```$")

// GenerateCode sends the provided prompt to the OpenAI API and returns the
// generated code.
func (client *Client) GenerateCode(ctx context.Context, prompt string) (
	code string,
	err error,
) {
	if client.model == ModelChatGPT {
		code, err = client.generateWithChatModel(ctx, prompt)
	} else {
		code, err = client.generateWithCompletionsModel(ctx, prompt)
	}

	if err != nil {
		return "", err
	}

	if client.full {
		return code, nil
	}

	m := codeRegex.FindStringSubmatch(code)
	if m == nil || m[1] == "" {
		return code, nil
	}

	return m[1], nil
}

func (client *Client) generateWithChatModel(ctx context.Context, prompt string) (
	code string,
	err error,
) {
	var answer struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			Index        int64  `json:"index"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}

	err = client.NewRequest("POST", "/chat/completions").
		JSONBody(map[string]interface{}{
			"model": client.model,
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
			"max_tokens": MaxTokens + 1 - len(prompt),
		}).
		Into(&answer).
		RunContext(ctx)
	if err != nil {
		return code, fmt.Errorf("failed sending prompt: %w", err)
	}

	if len(answer.Choices) == 0 {
		return code, ErrNoResults
	}

	if answer.Choices[0].FinishReason != "stop" {
		return code, fmt.Errorf(
			"%w: %s",
			ErrResultTruncated,
			answer.Choices[0].FinishReason,
		)
	}

	return strings.TrimSpace(answer.Choices[0].Message.Content), nil
}

func (client *Client) generateWithCompletionsModel(
	ctx context.Context,
	prompt string,
) (code string, err error) {
	var answer struct {
		Choices []struct {
			Text         string `json:"text"`
			Index        int64  `json:"index"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}

	err = client.NewRequest("POST", "/completions").
		JSONBody(map[string]interface{}{
			"model":      client.model,
			"prompt":     prompt,
			"max_tokens": MaxTokens + 1 - len(prompt),
		}).
		Into(&answer).
		RunContext(ctx)
	if err != nil {
		return code, fmt.Errorf("failed sending prompt: %w", err)
	}

	if len(answer.Choices) == 0 {
		return code, ErrNoResults
	}

	if answer.Choices[0].FinishReason != "stop" {
		return code, fmt.Errorf(
			"%w: %s",
			ErrResultTruncated,
			answer.Choices[0].FinishReason,
		)
	}

	return strings.TrimSpace(answer.Choices[0].Text), nil
}
