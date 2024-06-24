package ollama

import (
	"context"
	"fmt"
	"sort"

	"github.com/gofireflyio/aiac/v5/libaiac/types"
)

// ListModels returns a list of all the models supported by this backend.
func (backend *Ollama) ListModels(ctx context.Context) (models []string, err error) {
	var answer struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	err = backend.
		NewRequest("GET", "/tags").
		Into(&answer).
		RunContext(ctx)
	if err != nil {
		return models, fmt.Errorf("failed sending prompt: %w", err)
	}

	if len(answer.Models) == 0 {
		return models, types.ErrNoResults
	}

	models = make([]string, len(answer.Models))
	for i := range answer.Models {
		models[i] = answer.Models[i].Name
	}

	sort.Strings(models)

	return models, nil
}
