package types

import (
	"context"
)

// Backend is an interface that must be implemented in order to support an LLM
// provider.
type Backend interface {
	// ListModels returns a list of all models supported by the backend.
	ListModels(context.Context) ([]string, error)

	// Chat initiates a conversation with an LLM backend. The name of the model
	// to use must be provided. Users can also supply zero or more "previous
	// messages" that may have been exchanged in the past. This practically
	// allows "loading" previous conversations and continuing them.
	Chat(string, ...Message) Conversation
}

// Conversation is an interface that must be implemented in order to support
// chat models in an LLM provider.
type Conversation interface {
	// Send sends a message to the model and returns the response.
	Send(context.Context, string) (Response, error)

	// Messages returns all the messages that have been exchanged between the
	// user and the assistant up to this point.
	Messages() []Message

	// AddHeader adds an extra HTTP header that will be added to every HTTP
	// request issued as part of this conversation. Any headers added will be in
	// addition to any extra headers defined for the backend itself, and will
	// take precedence over them. Not all providers may support this
	// (specifically, bedrock doesn't).
	AddHeader(string, string)
}
