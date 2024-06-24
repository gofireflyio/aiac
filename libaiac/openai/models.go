package openai

import (
	"context"
	"fmt"
	"sort"

	"github.com/gofireflyio/aiac/v5/libaiac/types"
)

// ListModels returns a list of all the models supported by this backend.
func (backend *OpenAI) ListModels(ctx context.Context) (
	models []string,
	err error,
) {
	var answer struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	err = backend.
		NewRequest("GET", "/models").
		Into(&answer).
		RunContext(ctx)
	if err != nil {
		return models, fmt.Errorf("failed sending prompt: %w", err)
	}

	if len(answer.Data) == 0 {
		return models, types.ErrNoResults
	}

	models = make([]string, len(answer.Data))
	for i := range answer.Data {
		models[i] = answer.Data[i].ID
	}

	sort.Strings(models)

	return models, nil
}
