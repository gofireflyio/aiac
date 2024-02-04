package ollama

import (
	"context"
	"fmt"

	"github.com/gofireflyio/aiac/v4/libaiac/types"
	ol "github.com/jmorganca/ollama/api"
)

// Complete sends a prompt to the generate endpoint.
func (c *Client) Complete(ctx context.Context, model types.Model, prompt string) (types.Response, error) {
	res := types.Response{}

	handleResponse := func(resp ol.GenerateResponse) error {
		if !resp.Done {
			return types.ErrResultTruncated
		}

		res.FullOutput = resp.Response

		if code, ok := types.ExtractCode(res.FullOutput); ok {
			res.Code = code
		} else {
			res.Code = res.FullOutput
		}

		res.TokensUsed = int64(resp.EvalCount)

		return nil
	}

	stream := false
	req := &ol.GenerateRequest{
		Model:  model.Name,
		Prompt: prompt,
		Stream: &stream,
	}

	err := c.backend.Generate(ctx, req, handleResponse)
	if err != nil {
		return types.Response{}, fmt.Errorf("ollama request failed: %w", err)
	}

	return res, nil
}
