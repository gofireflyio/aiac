package ollama

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofireflyio/aiac/v4/libaiac/types"
)

type completionResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Complete sends a request to OpenAI's Completions API using the provided model
// and prompt, and returns the response
func (client *Client) Complete(
	ctx context.Context,
	model types.Model,
	prompt string,
) (res types.Response, err error) {
	var answer completionResponse

	err = client.NewRequest("POST", "/generate").
		JSONBody(map[string]interface{}{
			"model":  model.Name,
			"prompt": prompt,
			"options": map[string]interface{}{
				"temperature": 0.2,
			},
			"stream": false,
		}).
		Into(&answer).
		RunContext(ctx)
	if err != nil {
		return res, fmt.Errorf("failed sending prompt: %w", err)
	}

	if !answer.Done {
		return res, fmt.Errorf("%w: unexpected truncated response", types.ErrResultTruncated)
	}

	res.FullOutput = strings.TrimSpace(answer.Response)

	var ok bool
	if res.Code, ok = types.ExtractCode(res.FullOutput); !ok {
		res.Code = res.FullOutput
	}

	return res, nil
}
