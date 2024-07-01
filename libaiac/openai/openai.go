package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofireflyio/aiac/v5/libaiac/types"
	"github.com/ido50/requests"
)

// OpenAIBackend is the default URI endpoint for the OpenAI API
const OpenAIBackend = "https://api.openai.com/v1"

// OpenAI is a structure used to continuously generate IaC code via OpenAPI
type OpenAI struct {
	*requests.HTTPClient
	apiKey     string
	apiVersion string
	authHeader string
}

// Options is a struct containing all the parameters accepted by the New
// constructor.
type Options struct {
	// APIKey is the OpenAI API key to use for requests. This is required.
	ApiKey string

	// URL is the OpenAI API URL to userequests. Optional, defaults to OpenAIBackend.
	URL string

	// APIVersion is the version of the OpenAI API to use. Optional.
	APIVersion string

	// AuthHeader allows modifying the header where the API key is sent. This
	// defaults to Authorization. If it is "Authorization" or
	// "Proxy-Authorization", the API key is sent with a "Bearer " prefix. If
	// it's anything else, the API key is sent alone.
	AuthHeader string

	// ExtraHeaders are extra HTTP headers to send with every request to the
	// provider.
	ExtraHeaders map[string]string
}

// New creates a new instance of the OpenAI struct, with the provided input
// options. The OpenAI API backend is not yet contacted at this point.
func New(opts *Options) (*OpenAI, error) {
	if opts == nil {
		return nil, nil
	}

	if opts.ApiKey == "" {
		return nil, errors.New("OpenAI backends require an API key")
	}

	// Trim "Bearer " prefix if user accidentally included it, probably by
	// copy-pasting from somewhere.
	opts.ApiKey = strings.TrimPrefix(opts.ApiKey, "Bearer ")

	if opts.URL == "" {
		opts.URL = OpenAIBackend
	}

	authHeaderKey := "Authorization"
	authHeaderVal := fmt.Sprintf("Bearer %s", opts.ApiKey)

	// If user provided a different authorization header, use it, and if that
	// header is neither "Authorization" nor "Proxy-Authorization", remove the
	// "Bearer " prefix from its value.
	if opts.AuthHeader != "" && opts.AuthHeader != authHeaderKey {
		authHeaderKey = opts.AuthHeader
		if authHeaderKey != "Proxy-Authorization" {
			authHeaderVal = opts.ApiKey
		}
	}

	// The above section depends on the user telling us to use a different
	// header for authorization. Previously, though, we used 'api-key' as the
	// header if the URL was anything other than the OpenAI URL. This worked for
	// Azure OpenAI users, but since many more providers now implement the same
	// API (e.g. Portkey), that check was no longer correct. To maintain
	// backwards compatibility for Azure OpenAI users, though, we can change the
	// auth header by ourselves if the URL is *.openai.azure.com
	if opts.AuthHeader == "" && strings.Contains(opts.URL, ".openai.azure.com") {
		authHeaderKey = "api-key"
		authHeaderVal = opts.ApiKey
	}

	backend := &OpenAI{
		apiKey:     opts.ApiKey,
		apiVersion: opts.APIVersion,
	}

	backend.HTTPClient = requests.NewClient(opts.URL).
		Accept("application/json").
		Header(authHeaderKey, authHeaderVal).
		ErrorHandler(func(
			httpStatus int,
			contentType string,
			body io.Reader,
		) error {
			var res struct {
				Error struct {
					Message string `json:"Message"`
					Type    string `json:"type"`
				} `json:"error"`
			}

			err := json.NewDecoder(body).Decode(&res)
			if err != nil {
				return fmt.Errorf(
					"%w %s",
					types.ErrUnexpectedStatus,
					http.StatusText(httpStatus),
				)
			}

			return fmt.Errorf(
				"%w: [%s]: %s",
				types.ErrRequestFailed,
				res.Error.Type,
				res.Error.Message,
			)
		})

	for header, value := range opts.ExtraHeaders {
		backend.HTTPClient.Header(header, value)
	}

	return backend, nil
}
