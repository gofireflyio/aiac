package bedrock

import (
	"github.com/gofireflyio/aiac/v4/libaiac/types"
)

var (
	// ModelTitanG1Lite is the amazon.titan-text-lite-v1 model.
	ModelTitanG1Lite = types.Model{"amazon.titan-text-lite-v1", 4096, types.ModelTypeCompletion}

	// ModelTitanG1Express is the amazon.titan-text-express-v1 model.
	ModelTitanG1Express = types.Model{"amazon.titan-text-express-v1", 8192, types.ModelTypeChat}

	// ModelClaude1 is the anthropic.claude-v1 model.
	ModelClaude1 = types.Model{"anthropic.claude-v1", 100000, types.ModelTypeChat}

	// ModelClaude1 is the anthropic.claude-v2 model.
	ModelClaude2 = types.Model{"anthropic.claude-v2", 100000, types.ModelTypeChat}

	// ModelClaude1 is the anthropic.claude-v2:1 model.
	ModelClaude21 = types.Model{"anthropic.claude-v2:1", 200000, types.ModelTypeChat}

	// SupportedModels is a list of all language models supported by this
	// backend implementation.
	SupportedModels = []types.Model{
		ModelTitanG1Lite,
		ModelTitanG1Express,
		ModelClaude1,
		ModelClaude2,
		ModelClaude21,
	}
)

// ListModels returns a list of all the models supported by this backend
// implementation.
func (client *Client) ListModels() []types.Model {
	return SupportedModels
}

// DefaultModel returns the default model used by this backend.
func (client *Client) DefaultModel() types.Model {
	return ModelTitanG1Lite
}
