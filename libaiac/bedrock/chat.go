package bedrock

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gofireflyio/aiac/v4/libaiac/types"
)

// Conversation is a struct used to converse with a Bedrock chat model. It
// maintains all messages sent/received in order to maintain context.
type Conversation struct {
	client   *Client
	model    types.Model
	messages []types.Message
}

// Chat initiates a conversation with a Bedrock chat model. A conversation
// maintains context, allowing to send further instructions to modify the output
// from previous requests.
func (client *Client) Chat(model types.Model) types.Conversation {
	if model.Type != types.ModelTypeChat {
		return nil
	}

	return &Conversation{
		client: client,
		model:  model,
	}
}

// Send sends the provided message to the backend and returns a Response object.
// To maintain context, all previous messages (whether from you to the API or
// vice-versa) are sent as well, allowing you to ask the API to modify the
// code it already generated.
func (conv *Conversation) Send(ctx context.Context, prompt string, msgs ...types.Message) (
	res types.Response,
	err error,
) {
	if len(msgs) > 0 {
		conv.messages = append(conv.messages, msgs...)
	}

	conv.messages = append(conv.messages, types.Message{
		Role:    "user",
		Content: prompt,
	})

	var inputText strings.Builder
	for _, msg := range conv.messages {
		switch msg.Role {
		case "user":
			fmt.Fprint(&inputText, "\n\nHuman: ")
		default:
			fmt.Fprint(&inputText, "\n\nAssistant: ")
		}

		fmt.Fprint(&inputText, msg.Content)
	}

	fmt.Fprintf(&inputText, "\n\nAssistant:")

	body, err := conv.client.generateInputJSON(conv.model, inputText.String())
	if err != nil {
		return res, fmt.Errorf("failed generating input JSON: %w", err)
	}

	output, err := conv.client.backend.InvokeModel(
		ctx,
		&bedrockruntime.InvokeModelInput{
			Body:        body,
			ModelId:     aws.String(conv.model.Name),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		},
	)
	if err != nil {
		return res, fmt.Errorf("failed sending prompt: %w", err)
	}

	res.FullOutput, res.TokensUsed, err = conv.client.parseOutputJSON(
		conv.model,
		output.Body,
	)
	if err != nil {
		return res, err
	}

	conv.messages = append(conv.messages, types.Message{
		Role:    "assistant",
		Content: res.FullOutput,
	})

	var ok bool
	if res.Code, ok = types.ExtractCode(res.FullOutput); !ok {
		res.Code = res.FullOutput
	}

	return res, nil
}
