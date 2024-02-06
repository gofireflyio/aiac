package ollama

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofireflyio/aiac/v4/libaiac/types"
	"github.com/ido50/requests"
)

// DefaultAPIURL is the default URL for a local Ollama API server
const DefaultAPIURL = "http://localhost:11434/api"

// Client is a structure used to continuously generate IaC code via Ollama
type Client struct {
	*requests.HTTPClient
}

// NewClientOptions is a struct containing all the parameters accepted by the
// NewClient constructor.
type NewClientOptions struct {
	// URL is the URL of the API server (including the /api path prefix). Defaults to DefaultAPIURL.
	URL string
}

// NewClient creates a new instance of the Client struct, with the provided
// input options. The Ollama API server is not contacted at this point.
func NewClient(opts *NewClientOptions) *Client {
	if opts == nil {
		opts = &NewClientOptions{}
	}

	if opts.URL == "" {
		opts.URL = DefaultAPIURL
	}

	cli := &Client{}

	cli.HTTPClient = requests.NewClient(opts.URL).
		Accept("application/json").
		ErrorHandler(func(
			httpStatus int,
			contentType string,
			body io.Reader,
		) error {
			var res struct {
				Error string `json:"error"`
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
				"%w:  %s",
				types.ErrRequestFailed,
				res.Error,
			)
		})

	return cli
}
