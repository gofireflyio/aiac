package ollama

import (
	"context"
	"fmt"

	"github.com/gofireflyio/aiac/v4/libaiac/types"
	ol "github.com/jmorganca/ollama/api"
)

// Conversation is a struct used to converse with a chat model. It
// maintains all messages sent/received in order to maintain context.
type Conversation struct {
	client   *Client
	model    types.Model
	messages []ol.Message
}

// Chat initiates a conversation with a chat model.
//
//nolint:ireturn
func (c *Client) Chat(model types.Model) types.Conversation {
	if model.Type != types.ModelTypeChat {
		return nil
	}

	return &Conversation{
		client: c,
		model:  model,
	}
}

// Send sends a message to the model and returns the response.
func (c *Conversation) Send(ctx context.Context, prompt string, messages ...types.Message) (types.Response, error) {
	res := types.Response{}
	handleResponse := func(resp ol.ChatResponse) error {
		if !resp.Done {
			return types.ErrResultTruncated
		}

		res.FullOutput = resp.Message.Content

		if code, ok := types.ExtractCode(res.FullOutput); ok {
			res.Code = code
		} else {
			res.Code = res.FullOutput
		}

		res.TokensUsed = int64(resp.EvalCount)

		return nil
	}

	if len(messages) > 0 {
		for _, msg := range messages {
			c.messages = append(c.messages, ol.Message{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	c.messages = append(c.messages, ol.Message{
		Role:    "user",
		Content: prompt,
	})

	stream := false
	chatReq := &ol.ChatRequest{
		Model:    c.model.Name,
		Stream:   &stream,
		Messages: c.messages,
	}

	err := c.client.backend.Chat(ctx, chatReq, handleResponse)
	if err != nil {
		return types.Response{}, fmt.Errorf("failed sending chat request: %w", err)
	}

	c.messages = append(c.messages, ol.Message{
		Role:    "assistant",
		Content: res.FullOutput,
	})

	return res, nil
}
