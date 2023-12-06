package types

import "regexp"

// Message represents a single message in an exchange between a user and an
// AI model, either as part of a chat or a single completion request.
type Message struct {
	// Role is the type of the participant. The user is named "user" (in Amazon
	// Bedrock, this is equivalent to the "Human" identifier). Anything else is
	// considered the AI model.
	Role string `json:"role"`

	// Content is the text content of the message.
	Content string `json:"content"`
}

// Response is the struct returned from methods generating code via the OpenAI
// API.
type Response struct {
	// FullOutput is the complete output returned by the API. This is generally
	// a Markdown-formatted Message that contains the generated code, plus
	// explanations, if any.
	FullOutput string

	// Code is the extracted code section from the complete output. If code was
	// not found or extraction otherwise failed, this will be the same as
	// FullOutput.
	Code string

	// APIKeyUsed is the API key used when making the request.
	APIKeyUsed string

	// TokensUsed is the number of tokens utilized by the request. This is
	// the "usage.total_tokens" value returned from the API.
	TokensUsed int64
}

var codeRegex = regexp.MustCompile("(?ms)^```(?:[^\n]*)\n(.*?)\n```$")

// ExtractCode receives the full output string from the OpenAI API and attempts
// to extract a code block from it. OpenAI code blocks are generally Markdown
// blocks surrounded by the ``` string on both sides. If successful, the code
// string will be returned together with a true value, otherwise an empty string
// is returned together with a false value.
func ExtractCode(output string) (string, bool) {
	m := codeRegex.FindStringSubmatch(output)
	if m == nil || m[1] == "" {
		return "", false
	}

	return m[1], true
}
