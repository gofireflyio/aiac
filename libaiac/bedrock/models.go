package bedrock

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrock/types"
)

// ListModels returns a list of all the models supported by this backend.
func (backend *Bedrock) ListModels(ctx context.Context) (models []string, err error) {
	output, err := backend.service.ListFoundationModels(ctx, &bedrock.ListFoundationModelsInput{
		ByOutputModality: types.ModelModalityText,
	})
	if err != nil {
		return models, fmt.Errorf("failed listing base models: %w", err)
	}

	models = make([]string, len(output.ModelSummaries))
	for i := range output.ModelSummaries {
		models[i] = *output.ModelSummaries[i].ModelId
	}

	sort.Strings(models)

	return models, nil
}
