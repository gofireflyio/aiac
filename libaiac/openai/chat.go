package openai

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofireflyio/aiac/v4/libaiac/types"
)

// Conversation is a struct used to converse with an OpenAI chat model. It
// maintains all messages sent/received in order to maintain context just like
// using ChatGPT.
type Conversation struct {
	client   *Client
	model    types.Model
	messages []types.Message
}

type chatResponse struct {
	Choices []struct {
		Message      types.Message `json:"message"`
		Index        int64         `json:"index"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int64 `json:"total_tokens"`
	} `json:"usage"`
}

// Chat initiates a conversation with an OpenAI chat model. A conversation
// maintains context, allowing to send further instructions to modify the output
// from previous requests, just like using the ChatGPT website.
func (client *Client) Chat(model types.Model) types.Conversation {
	if model.Type != types.ModelTypeChat {
		return nil
	}

	return &Conversation{
		client: client,
		model:  model,
	}
}

// Send sends the provided message to the API and returns a Response object.
// To maintain context, all previous messages (whether from you to the API or
// vice-versa) are sent as well, allowing you to ask the API to modify the
// code it already generated.
func (conv *Conversation) Send(ctx context.Context, prompt string, msgs ...types.Message) (
	res types.Response,
	err error,
) {
	var answer chatResponse

	if len(msgs) > 0 {
		conv.messages = append(conv.messages, msgs...)
	}

	conv.messages = append(conv.messages, types.Message{
		Role:    "user",
		Content: prompt,
	})

	var apiVersion string
	if len(conv.client.apiVersion) > 0 {
		apiVersion = fmt.Sprintf("?api-version=%s", conv.client.apiVersion)
	}

	err = conv.client.NewRequest("POST",
		fmt.Sprintf("/chat/completions%s", apiVersion)).
		JSONBody(map[string]interface{}{
			"model":       conv.model.Name,
			"messages":    conv.messages,
			"temperature": 0.2,
		}).
		Into(&answer).
		RunContext(ctx)
	if err != nil {
		return res, fmt.Errorf("failed sending prompt: %w", err)
	}

	if len(answer.Choices) == 0 {
		return res, types.ErrNoResults
	}

	if answer.Choices[0].FinishReason != "stop" {
		return res, fmt.Errorf(
			"%w: %s",
			types.ErrResultTruncated,
			answer.Choices[0].FinishReason,
		)
	}

	conv.messages = append(conv.messages, answer.Choices[0].Message)

	res.FullOutput = strings.TrimSpace(answer.Choices[0].Message.Content)
	res.APIKeyUsed = conv.client.apiKey
	res.TokensUsed = answer.Usage.TotalTokens

	var ok bool
	if res.Code, ok = types.ExtractCode(res.FullOutput); !ok {
		res.Code = res.FullOutput
	}

	return res, nil
}
