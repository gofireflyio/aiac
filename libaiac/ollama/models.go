package ollama

import (
	"context"
	"time"

	"github.com/gofireflyio/aiac/v4/libaiac/types"
	ol "github.com/jmorganca/ollama/api"
)

var defaultOllamaContextSize = ol.DefaultOptions().NumCtx
var timeout = 5 * time.Second

func (c *Client) ListModels() []types.Model {
	if c.backend == nil {
		c.backend = newBackend()
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	lr, err := c.backend.List(ctx)

	if err != nil {
		panic(err)
	}

	models := []types.Model{}

	for _, model := range lr.Models {
		req := &ol.ShowRequest{
			Model: model.Model,
		}

		resp, err := c.backend.Show(ctx, req)
		if err != nil {
			break
		}

		// ollama defaults to 2048 token, but can be configured
		// we could take max tokens as a  configuration param
		models = append(models, types.Model{
			Name:      model.Model,
			MaxTokens: defaultOllamaContextSize,
			Type:      getModelType(resp),
		})
	}

	return models
}

func getModelType(showResponse *ol.ShowResponse) types.ModelType {
	// Lets make a basic assumption that if the model template just
	// takes the prompt with nothing else then it is completion,
	// otherwise it is chat.
	if showResponse.Template == "{{ .Prompt }}" {
		return types.ModelTypeCompletion
	}

	return types.ModelTypeChat
}

func (c *Client) DefaultModel() types.Model {
	if c.backend == nil {
		c.backend = newBackend()
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	lr, err := c.backend.List(ctx)

	if err != nil {
		panic(err)
	}

	for _, model := range lr.Models {
		req := &ol.ShowRequest{
			Model: model.Model,
		}

		resp, err := c.backend.Show(ctx, req)

		if err != nil {
			continue
		}

		return types.Model{
			Name:      model.Model,
			MaxTokens: defaultOllamaContextSize,
			Type:      getModelType(resp),
		}
	}

	return types.Model{}
}
