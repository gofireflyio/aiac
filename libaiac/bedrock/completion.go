package bedrock

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gofireflyio/aiac/v4/libaiac/types"
)

func (client *Client) generateInputJSON(
	model types.Model,
	prompt string,
) (body []byte, err error) {
	switch model {
	case ModelTitanG1Lite, ModelTitanG1Express:
		body, err = json.Marshal(map[string]interface{}{
			"inputText": prompt,
			"textGenerationConfig": map[string]interface{}{
				"temperature":   0.2,
				"maxTokenCount": model.MaxTokens,
			},
		})
	case ModelClaude1, ModelClaude2, ModelClaude21:
		body, err = json.Marshal(map[string]interface{}{
			"prompt":               prompt,
			"temperature":          0.2,
			"max_tokens_to_sample": 4000, // as currently recommended by Amazon
			"stop_sequences":       []string{"\n\nHuman:"},
		})
		if err != nil {
			return body, fmt.Errorf("failed generating request JSON: %w", err)
		}
	default:
		return body, types.ErrUnsupportedModel
	}

	return body, err
}

func (client *Client) parseOutputJSON(model types.Model, body []byte) (
	text string,
	tokensUsed int64,
	err error,
) {
	switch model {
	case ModelTitanG1Lite, ModelTitanG1Express:
		var bodyStruct struct {
			InputTokenCount int64 `json:"inputTextTokenCount"`
			Results         []struct {
				TokenCount       int64  `json:"tokenCount"`
				OutputText       string `json:"outputText"`
				CompletionReason string `json:"completionReason"`
			} `json:"results"`
		}

		err = json.Unmarshal(body, &bodyStruct)
		if err != nil {
			return text, tokensUsed, fmt.Errorf("failed decoding response: %w", err)
		}

		if len(bodyStruct.Results) == 0 {
			return text, tokensUsed, types.ErrNoResults
		}

		if bodyStruct.Results[0].CompletionReason != "FINISH" {
			return text, tokensUsed, fmt.Errorf(
				"%w: %s",
				types.ErrResultTruncated,
				bodyStruct.Results[0].CompletionReason,
			)
		}

		tokensUsed = bodyStruct.InputTokenCount
		for _, choice := range bodyStruct.Results {
			tokensUsed += choice.TokenCount
		}

		text = strings.TrimSpace(bodyStruct.Results[0].OutputText)
	case ModelClaude1, ModelClaude2, ModelClaude21:
		var bodyStruct struct {
			OutputText       string `json:"completion"`
			CompletionReason string `json:"stop_reason"`
		}

		err = json.Unmarshal(body, &bodyStruct)
		if err != nil {
			return text, tokensUsed, fmt.Errorf("failed decoding response: %w", err)
		}

		if bodyStruct.CompletionReason != "stop_sequence" {
			return text, tokensUsed, fmt.Errorf(
				"%w: %s",
				types.ErrResultTruncated,
				bodyStruct.CompletionReason,
			)
		}

		return bodyStruct.OutputText, 0, nil
	}

	return text, tokensUsed, nil
}

// Complete sends a request to a Bedrock completion model with the provided
// prompt, and returns the response
func (client *Client) Complete(
	ctx context.Context,
	model types.Model,
	prompt string,
) (res types.Response, err error) {
	body, err := client.generateInputJSON(model, prompt)
	if err != nil {
		return res, fmt.Errorf("failed generating input JSON: %w", err)
	}

	output, err := client.backend.InvokeModel(
		ctx,
		&bedrockruntime.InvokeModelInput{
			Body:        body,
			ModelId:     aws.String(model.Name),
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		},
	)
	if err != nil {
		return res, fmt.Errorf("failed sending prompt: %w", err)
	}

	res.FullOutput, res.TokensUsed, err = client.parseOutputJSON(model, output.Body)
	if err != nil {
		return res, err
	}

	var ok bool
	if res.Code, ok = types.ExtractCode(res.FullOutput); !ok {
		res.Code = res.FullOutput
	}

	return res, nil
}
