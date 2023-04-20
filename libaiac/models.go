package libaiac

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

var (
	// ModelGPT35Turbo represents the gpt-3.5-turbo model used by ChatGPT
	ModelGPT35Turbo = Model{"gpt-3.5-turbo", 4096, ModelTypeChat}

	// ModelGPT35Turbo represents the gpt-3.5-turbo-0301 model, a March 1st 2023
	// snapshot of gpt-3.5-turbo
	ModelGPT35Turbo0301 = Model{"gpt-3.5-turbo-0301", 4096, ModelTypeChat}

	// ModelGPT4 represents the gpt-4 model
	ModelGPT4 = Model{"gpt-4", 8192, ModelTypeChat}

	// ModelGPT40314 represents the gpt-4-0314 model, a March 14th 2023 snapshot
	// of the gpt-4 model.
	ModelGPT40314 = Model{"gpt-4-0314", 8192, ModelTypeChat}

	// ModelGPT432K represents the gpt-4-32k model, which is the same as gpt-4,
	// but with 4x the context length.
	ModelGPT432K = Model{"gpt-4-32k", 32768, ModelTypeChat}

	// ModelGPT432K0314 represents the gpt-4-32k-0314 model, a March 14th 2023
	// snapshot of the gpt-4-32k model
	ModelGPT432K0314 = Model{"gpt-4-32k-0314", 32768, ModelTypeChat}

	// ModelTextDaVinci3 represents the text-davinci-003 language generation
	// model.
	ModelTextDaVinci3 = Model{"text-davinci-003", 4097, ModelTypeCompletion}

	// ModelTextDaVinci2 represents the text-davinci-002 language generation
	// model.
	ModelTextDaVinci2 = Model{"text-davinci-002", 4097, ModelTypeCompletion}

	// SupportedModels is a list of all language models supported by aiac
	SupportedModels = []Model{
		ModelGPT35Turbo,
		ModelGPT35Turbo0301,
		ModelGPT4,
		ModelGPT40314,
		ModelGPT432K,
		ModelGPT432K0314,
		ModelTextDaVinci3,
		ModelTextDaVinci2,
	}
)

// Decode is used by the kong library to map CLI-provided values to the Model
// type
func (m *Model) Decode(ctx *kong.DecodeContext) error {
	var provided string

	err := ctx.Scan.PopValueInto("string", &provided)
	if err != nil {
		return fmt.Errorf("failed getting model value: %w", err)
	}

	if provided == "" {
		*m = ModelGPT35Turbo
		return nil
	}

	for _, supported := range SupportedModels {
		if supported.Name == provided {
			*m = supported
			return nil
		}
	}

	return fmt.Errorf("%w %s", ErrUnsupportedModel, provided)
}
