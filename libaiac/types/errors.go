package types

import "errors"

var (
	// ErrResultTruncated is returned when the OpenAI API returned a truncated
	// result. The reason for the truncation will be appended to the error
	// string.
	ErrResultTruncated = errors.New("result was truncated")

	// ErrNoResults is returned if the OpenAI API returned an empty result. This
	// should not generally happen.
	ErrNoResults = errors.New("no results return from API")

	// ErrUnsupportedBackend is returned if the provided backend name is
	// unknown.
	ErrUnsupportedBackend = errors.New("unsupported backend")

	// ErrUnsupportedModel is returned if the SetModel method is provided with
	// an unsupported model
	ErrUnsupportedModel = errors.New("unsupported model")

	// ErrUnexpectedStatus is returned when the OpenAI API returned a response
	// with an unexpected status code
	ErrUnexpectedStatus = errors.New("OpenAI returned unexpected response")

	// ErrRequestFailed is returned when the OpenAI API returned an error for
	// the request
	ErrRequestFailed = errors.New("request failed")
)
