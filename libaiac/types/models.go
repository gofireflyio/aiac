package types

import (
	"fmt"

	"github.com/alecthomas/kong"
)

// ModelType is an enum indicating the type of language models
type ModelType string

const (
	// ModelTypeCompletion is used to represent text completion models
	ModelTypeCompletion ModelType = "completion"

	// ModelTypeChat is used to represent chat models
	ModelTypeChat ModelType = "chat"
)

// Decode is used by the kong library to map CLI-provided values to the Model
// type
func (m *ModelType) Decode(ctx *kong.DecodeContext) error {
	var provided string

	err := ctx.Scan.PopValueInto("string", &provided)
	if err != nil {
		return fmt.Errorf("failed getting model value: %w", err)
	}

	switch provided {
	case string(ModelTypeCompletion):
		*m = ModelTypeCompletion
	case string(ModelTypeChat):
		*m = ModelTypeChat
	default:
		return fmt.Errorf("%w %s", ErrUnsupportedModel, provided)
	}

	return nil
}

// Model is a struct used to represent supported language models
type Model struct {
	Name      string
	MaxTokens int
	Type      ModelType
}
