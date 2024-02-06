package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofireflyio/aiac/v4/libaiac/types"
	"github.com/ido50/requests"
)

// OpenAIBackend is the default URI endpoint for the OpenAI API.
const OpenAIBackend = "https://api.openai.com/v1"

// Client is a structure used to continuously generate IaC code via OpenAPI
type Client struct {
	*requests.HTTPClient
	apiKey     string
	apiVersion string
}

// NewClientOptions is a struct containing all the parameters accepted by the
// NewClient constructor.
type NewClientOptions struct {
	// APIKey is the OpenAI API key to use for requests. This is required.
	ApiKey string

	// URL is the OpenAI API URL to userequests. Optional, defaults to OpenAIBackend.
	URL string

	// APIVersion is the version of the OpenAI API to use. Optional.
	APIVersion string
}

// NewClient creates a new instance of the Client struct, with the provided
// input options. The OpenAI API backend is not yet contacted at this point.
func NewClient(opts *NewClientOptions) *Client {
	if opts == nil {
		return nil
	}

	if opts.ApiKey == "" {
		return nil
	}

	if opts.URL == "" {
		opts.URL = OpenAIBackend
	}

	var authHeaderKey string
	var authHeaderVal string

	if opts.URL == OpenAIBackend {
		authHeaderKey = "Authorization"
		authHeaderVal = fmt.Sprintf("Bearer %s", opts.ApiKey)
	} else {
		authHeaderKey = "api-key"
		authHeaderVal = opts.ApiKey
	}

	cli := &Client{
		apiKey:     strings.TrimPrefix(opts.ApiKey, "Bearer "),
		apiVersion: opts.APIVersion,
	}

	cli.HTTPClient = requests.NewClient(opts.URL).
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

	return cli
}
