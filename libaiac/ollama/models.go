package ollama

import (
	"github.com/gofireflyio/aiac/v4/libaiac/types"
)

var (
	// ModelCodeLlama represents the codellama model
	ModelCodeLlama = types.Model{"codellama", 0, types.ModelTypeChat}

	// ModelDeepseekCoder represents the deepseek-coder model
	ModelDeepseekCoder = types.Model{"deepseek-coder", 0, types.ModelTypeChat}

	// ModelWizardCoder represents the wizard-coder model
	ModelWizardCoder = types.Model{"wizard-coder", 0, types.ModelTypeChat}

	// ModelPhindCodeLlama represents the phind-codellama model
	ModelPhindCodeLlama = types.Model{"phind-codellama", 0, types.ModelTypeChat}

	// ModeCodeUp represents the codeup model
	ModelCodeUp = types.Model{"codeup", 0, types.ModelTypeChat}

	// ModeStarCoder represents the starcoder model
	ModelStarCoder = types.Model{"starcoder", 0, types.ModelTypeChat}

	// ModelSQLCoder represents the sqlcoder model
	ModelSQLCoder = types.Model{"sqlcoder", 0, types.ModelTypeChat}

	// ModelStableCode represents the stablecode model
	ModelStableCode = types.Model{"stablecode", 0, types.ModelTypeChat}

	// ModelMagicoder represents the magicoder model
	ModelMagicoder = types.Model{"magicoder", 0, types.ModelTypeChat}

	// ModelCodeBooga represents the codebooga model
	ModelCodeBooga = types.Model{"codebooga", 0, types.ModelTypeChat}

	// ModelMistral represents the mistral model
	ModelMistral = types.Model{"mistral", 0, types.ModelTypeCompletion}

	// SupportedModels is a list of all language models supported by this
	// backend implementation.
	SupportedModels = []types.Model{
		ModelCodeLlama,
		ModelDeepseekCoder,
		ModelWizardCoder,
		ModelPhindCodeLlama,
		ModelCodeUp,
		ModelStarCoder,
		ModelSQLCoder,
		ModelStableCode,
		ModelMagicoder,
		ModelCodeBooga,
		ModelMistral,
	}
)

// ListModels returns a list of all the models supported by this backend
// implementation.
func (client *Client) ListModels() []types.Model {
	return SupportedModels
}

// DefaultModel returns the default model used by this backend.
func (client *Client) DefaultModel() types.Model {
	return ModelMistral
}
