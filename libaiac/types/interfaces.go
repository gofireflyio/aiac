package types

import "context"

// Backend is an interface that must be implemented in order to support an LLM
// provider.
type Backend interface {
	// ListModels returns a list of all models supported by the backend.
	ListModels() []Model

	// DefaultModel returns the default model that should be used in the abscence
	// of a specific choice by the user.
	DefaultModel() Model

	// Complete send a prompt to a completion model.
	Complete(context.Context, Model, string) (Response, error)

	// Chat initiates a conversation with a chat model.
	Chat(Model) Conversation
}

// Conversation is an interface that must be implemented in order to support
// chat models in an LLM provider.
type Conversation interface {
	// Send sends a message to the model and returns the response.
	Send(context.Context, string, ...Message) (Response, error)
}
