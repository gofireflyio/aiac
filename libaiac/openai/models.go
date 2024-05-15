package openai

import (
	"github.com/gofireflyio/aiac/v4/libaiac/types"
)

var (
	// ModelGPT35Turbo represents the gpt-3.5-turbo model used by ChatGPT
	ModelGPT35Turbo = types.Model{"gpt-3.5-turbo", 4096, types.ModelTypeChat}

	// ModelGPT35Turbo represents the gpt-3.5-turbo-0301 model, a March 1st 2023
	// snapshot of gpt-3.5-turbo
	ModelGPT35Turbo0301 = types.Model{"gpt-3.5-turbo-0301", 4096, types.ModelTypeChat}

	// ModelGPT4 represents the gpt-4 model
	ModelGPT4 = types.Model{"gpt-4", 8192, types.ModelTypeChat}

	// ModelGPT40314 represents the gpt-4-0314 model, a March 14th 2023 snapshot
	// of the gpt-4 model.
	ModelGPT40314 = types.Model{"gpt-4-0314", 8192, types.ModelTypeChat}

	// ModelGPT432K represents the gpt-4-32k model, which is the same as gpt-4,
	// but with 4x the context length.
	ModelGPT432K = types.Model{"gpt-4-32k", 32768, types.ModelTypeChat}

	// ModelGPT432K0314 represents the gpt-4-32k-0314 model, a March 14th 2023
	// snapshot of the gpt-4-32k model
	ModelGPT432K0314 = types.Model{"gpt-4-32k-0314", 32768, types.ModelTypeChat}

	// ModelGPT4o represents the gpt-4o model
	ModelGPT4o = types.Model{"gpt-4o", 128000, types.ModelTypeChat}

	// ModelGPT4o20240513 represents the gpt-4o-2024-05-13 model
	ModelGPT4o20240513 = types.Model{"gpt-4o-2024-05-13", 128000, types.ModelTypeChat}

	// ModelGPT4Turbo represents the gpt-4-turbo model
	ModelGPT4Turbo = types.Model{"gpt-4-turbo", 128000, types.ModelTypeChat}

	// ModelGPT4Turbo20240409 represents the gpt-4-turbo-2024-04-09 model
	ModelGPT4Turbo20240409 = types.Model{"gpt-4-turbo-2024-04-09", 128000, types.ModelTypeChat}

	// SupportedModels is a list of all language models supported by this
	// backend implementation.
	SupportedModels = []types.Model{
		ModelGPT35Turbo,
		ModelGPT35Turbo0301,
		ModelGPT4,
		ModelGPT40314,
		ModelGPT432K,
		ModelGPT432K0314,
		ModelGPT4Turbo,
		ModelGPT4Turbo20240409,
		ModelGPT4o,
		ModelGPT4o20240513,
	}
)

// ListModels returns a list of all the models supported by this backend
// implementation.
func (client *Client) ListModels() []types.Model {
	return SupportedModels
}

// DefaultModel returns the default model used by this backend.
func (client *Client) DefaultModel() types.Model {
	return ModelGPT35Turbo
}
