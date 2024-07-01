package bedrock

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrocktypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/gofireflyio/aiac/v5/libaiac/types"
)

// Conversation is a struct used to converse with a Bedrock chat model. It
// maintains all messages sent/received in order to maintain context.
type Conversation struct {
	backend  *Bedrock
	model    string
	messages []bedrocktypes.Message
}

// Chat initiates a conversation with a Bedrock chat model. A conversation
// maintains context, allowing to send further instructions to modify the output
// from previous requests. The name of the model to use must be provided. Users
// can also supply zero or more "previous messages" that may have been exchanged
// in the past. This practically allows "loading" previous conversations and
// continuing them.
func (backend *Bedrock) Chat(model string, msgs ...types.Message) types.Conversation {
	conv := &Conversation{
		backend: backend,
		model:   model,
	}

	if len(msgs) > 0 {
		conv.messages = make([]bedrocktypes.Message, len(msgs))
		for i := range msgs {
			role := bedrocktypes.ConversationRoleUser
			if msgs[i].Role == "assistant" {
				role = bedrocktypes.ConversationRoleAssistant
			}

			conv.messages[i] = bedrocktypes.Message{
				Role: role,
				Content: []bedrocktypes.ContentBlock{
					&bedrocktypes.ContentBlockMemberText{Value: msgs[i].Content},
				},
			}
		}
	}

	return conv
}

// Send sends the provided message to the backend and returns a Response object.
// To maintain context, all previous messages (whether from you to the API or
// vice-versa) are sent as well, allowing you to ask the API to modify the
// code it already generated.
func (conv *Conversation) Send(ctx context.Context, prompt string) (
	res types.Response,
	err error,
) {
	conv.messages = append(conv.messages, bedrocktypes.Message{
		Role: bedrocktypes.ConversationRoleUser,
		Content: []bedrocktypes.ContentBlock{
			&bedrocktypes.ContentBlockMemberText{Value: prompt},
		},
	})

	input := bedrockruntime.ConverseInput{
		ModelId:  aws.String(conv.model),
		Messages: conv.messages,
		InferenceConfig: &bedrocktypes.InferenceConfiguration{
			Temperature: aws.Float32(0.2),
		},
	}

	output, err := conv.backend.runtime.Converse(ctx, &input)
	if err != nil {
		return res, fmt.Errorf("failed sending prompt: %w", err)
	}

	outputMsgMember, ok := output.Output.(*bedrocktypes.ConverseOutputMemberMessage)
	if !ok {
		return res, fmt.Errorf("Bedrock returned an unexpected response")
	}

	if len(outputMsgMember.Value.Content) == 0 {
		return res, fmt.Errorf("Bedrock didn't return any message")
	}

	outputMsg := outputMsgMember.Value

	outputTxt, ok := outputMsg.Content[0].(*bedrocktypes.ContentBlockMemberText)
	if !ok {
		return res, fmt.Errorf("Bedrock return an unexpected response")
	}

	res.FullOutput = outputTxt.Value
	res.TokensUsed = int64(*output.Usage.TotalTokens)
	res.StopReason = string(output.StopReason)

	conv.messages = append(conv.messages, outputMsg)

	if res.Code, ok = types.ExtractCode(res.FullOutput); !ok {
		res.Code = res.FullOutput
	}

	return res, nil
}

// Messages returns all the messages that have been exchanged between the user
// and the assistant up to this point.
func (conv *Conversation) Messages() []types.Message {
	msgs := make([]types.Message, len(conv.messages))
	for i, m := range conv.messages {
		content, _ := m.Content[0].(*bedrocktypes.ContentBlockMemberText)
		msgs[i] = types.Message{
			Role:    string(m.Role),
			Content: content.Value,
		}
	}
	return msgs
}
