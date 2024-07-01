package ollama

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofireflyio/aiac/v5/libaiac/types"
	"github.com/ido50/requests"
)

// DefaultAPIURL is the default URL for a local Ollama API server
const DefaultAPIURL = "http://localhost:11434/api"

// Ollama is a structure used to continuously generate IaC code via Ollama
type Ollama struct {
	*requests.HTTPClient
}

// Options is a struct containing all the parameters accepted by the New
// constructor.
type Options struct {
	// URL is the URL of the API server (including the /api path prefix).
	// Defaults to DefaultAPIURL.
	URL string

	// ExtraHeaders are extra HTTP headers to send with every request to the
	// provider.
	ExtraHeaders map[string]string
}

// New creates a new instance of the Ollama struct, with the provided
// input options. The Ollama API server is not contacted at this point.
func New(opts *Options) *Ollama {
	if opts == nil {
		opts = &Options{}
	}

	if opts.URL == "" {
		opts.URL = DefaultAPIURL
	}

	cli := &Ollama{}

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

	for header, value := range opts.ExtraHeaders {
		cli.HTTPClient.Header(header, value)
	}

	return cli
}
