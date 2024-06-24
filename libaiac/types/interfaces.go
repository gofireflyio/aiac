package types

import "context"

// Backend is an interface that must be implemented in order to support an LLM
// provider.
type Backend interface {
	// ListModels returns a list of all models supported by the backend.
	ListModels(context.Context) ([]string, error)

	// Chat initiates a conversation with an LLM backend.
	Chat(string) Conversation
}

// Conversation is an interface that must be implemented in order to support
// chat models in an LLM provider.
type Conversation interface {
	// Send sends a message to the model and returns the response.
	Send(context.Context, string) (Response, error)
}
