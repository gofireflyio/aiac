package libaiac

import (
	"context"
	"fmt"
	"strings"
)

type completionResponse struct {
	Choices []struct {
		Text         string `json:"text"`
		Index        int64  `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int64 `json:"total_tokens"`
	} `json:"usage"`
}

// Complete sends a request to OpenAI's Completions API using the provided model
// and prompt, and returns the response
func (client *Client) Complete(
	ctx context.Context,
	model Model,
	prompt string,
) (res Response, err error) {
	var answer completionResponse

	err = client.NewRequest("POST", "/completions").
		JSONBody(map[string]interface{}{
			"model":       model.Name,
			"prompt":      prompt,
			"max_tokens":  model.MaxTokens - len(prompt),
			"temperature": 0.2,
		}).
		Into(&answer).
		RunContext(ctx)
	if err != nil {
		return res, fmt.Errorf("failed sending prompt: %w", err)
	}

	if len(answer.Choices) == 0 {
		return res, ErrNoResults
	}

	if answer.Choices[0].FinishReason != "stop" {
		return res, fmt.Errorf(
			"%w: %s",
			ErrResultTruncated,
			answer.Choices[0].FinishReason,
		)
	}

	res.FullOutput = strings.TrimSpace(answer.Choices[0].Text)
	res.APIKeyUsed = client.apiKey
	res.TokensUsed = answer.Usage.TotalTokens

	var ok bool
	if res.Code, ok = ExtractCode(res.FullOutput); !ok {
		res.Code = res.FullOutput
	}

	return res, nil
}
