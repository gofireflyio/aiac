package openai

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofireflyio/aiac/v5/libaiac/types"
)

// Conversation is a struct used to converse with an OpenAI chat model. It
// maintains all messages sent/received in order to maintain context just like
// using ChatGPT.
type Conversation struct {
	// Messages is the list of all messages exchanged between the user and the
	// assistant.
	Messages []types.Message

	backend *OpenAI
	model   string
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
// from previous requests, just like using the ChatGPT website. The name of the
// model to use must be provided. Users can also supply zero or more "previous
// messages" that may have been exchanged in the past. This practically allows
// "loading" previous conversations and continuing them.
func (backend *OpenAI) Chat(model string, msgs ...types.Message) types.Conversation {
	chat := &Conversation{
		backend: backend,
		model:   model,
	}

	if len(msgs) > 0 {
		chat.Messages = msgs
	}

	return chat
}

// Send sends the provided message to the API and returns a Response object.
// To maintain context, all previous messages (whether from you to the API or
// vice-versa) are sent as well, allowing you to ask the API to modify the
// code it already generated.
func (conv *Conversation) Send(ctx context.Context, prompt string) (
	res types.Response,
	err error,
) {
	var answer chatResponse

	conv.Messages = append(conv.Messages, types.Message{
		Role:    "user",
		Content: prompt,
	})

	var apiVersion string
	if len(conv.backend.apiVersion) > 0 {
		apiVersion = fmt.Sprintf("?api-version=%s", conv.backend.apiVersion)
	}

	err = conv.backend.
		NewRequest("POST", fmt.Sprintf("/chat/completions%s", apiVersion)).
		JSONBody(map[string]interface{}{
			"model":       conv.model,
			"messages":    conv.Messages,
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

	conv.Messages = append(conv.Messages, answer.Choices[0].Message)

	res.FullOutput = strings.TrimSpace(answer.Choices[0].Message.Content)
	res.APIKeyUsed = conv.backend.apiKey
	res.TokensUsed = answer.Usage.TotalTokens
	res.StopReason = answer.Choices[0].FinishReason

	var ok bool
	if res.Code, ok = types.ExtractCode(res.FullOutput); !ok {
		res.Code = res.FullOutput
	}

	return res, nil
}
