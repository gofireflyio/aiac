package libaiac

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gofireflyio/aiac/v4/libaiac/bedrock"
	"github.com/gofireflyio/aiac/v4/libaiac/openai"
	"github.com/gofireflyio/aiac/v4/libaiac/types"
)

// Version contains aiac's version string
var Version = "development"

// Client provides the main interface for using libaiac. It exposes all the
// capabilities of the chosen backend.
type Client struct {
	// Backend is the backend implementation in charge of communicating with
	// the relevant APIs.
	Backend types.Backend
}

// BackendName is a const type used for identifying backends, a.k.a LLM providers.
type BackendName string

const (
	// BackendOpenAI represents the OpenAI LLM provider.
	BackendOpenAI BackendName = "openai"

	// BackendBedrock represents the Amazon Bedrock LLM provider.
	BackendBedrock BackendName = "bedrock"
)

// Decode is used by the kong library to map CLI-provided values to the Model
// type
func (b *BackendName) Decode(ctx *kong.DecodeContext) error {
	var provided string

	err := ctx.Scan.PopValueInto("string", &provided)
	if err != nil {
		return fmt.Errorf("failed getting model value: %w", err)
	}

	switch provided {
	case string(BackendOpenAI):
		*b = BackendOpenAI
	case string(BackendBedrock):
		*b = BackendBedrock
	default:
		return fmt.Errorf("%w %s", types.ErrUnsupportedBackend, provided)
	}

	return nil
}

// NewClientOptions contains all the parameters accepted by the NewClient
// constructor.
type NewClientOptions struct {
	// Backend is the name of the backend to use. Use the available constants,
	// e.g. BackendOpenAI or BackendBedrock. Defaults to openai.
	Backend BackendName

	// ----------------------
	// OpenAI related options
	// ----------------------

	// ApiKey is the OpenAI API key. Required if using OpenAI.
	ApiKey string

	// URL can be used to change the OpenAPI endpoint, for example in order to
	// use Azure OpenAI services. Defaults to OpenAI's standard API endpoint.
	URL string

	// APIVersion is the version of the OpenAI API to use. Unset by default.
	APIVersion string

	// ---------------------
	// Bedrock configuration
	// ---------------------

	// AWSRegion is the name of the region to use. Defaults to "us-east-1".
	AWSRegion string

	// AWSProfile is the name of the AWS profile to use. Defaults to "default".
	AWSProfile string
}

const (
	DefaultAWSRegion  = "us-east-1"
	DefaultAWSProfile = "default"
)

// NewClient constructs a new Client object.
func NewClient(opts *NewClientOptions) *Client {
	var backend types.Backend

	switch opts.Backend {
	case BackendBedrock:
		if opts.AWSProfile == "" {
			opts.AWSProfile = DefaultAWSProfile
		}
		if opts.AWSRegion == "" {
			opts.AWSRegion = DefaultAWSRegion
		}

		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithSharedConfigProfile(opts.AWSProfile),
		)
		if err != nil {
			return nil
		}

		backend = bedrock.NewClient(bedrockruntime.Options{
			Credentials: cfg.Credentials,
			Region:      opts.AWSRegion,
		})
	default:
		// default to openai
		backend = openai.NewClient(&openai.NewClientOptions{
			ApiKey:     opts.ApiKey,
			URL:        opts.URL,
			APIVersion: opts.APIVersion,
		})
	}

	return &Client{
		Backend: backend,
	}
}

// ListModels returns a list of all the models supported by the chosen backend
// implementation.
func (client *Client) ListModels() []types.Model {
	return client.Backend.ListModels()
}

// DefaultModel returns the default model used by the chosen backend implementation.
func (client *Client) DefaultModel() types.Model {
	return client.Backend.DefaultModel()
}

// Complete issues a request to a code completion model in the backend. with
// the provided string prompt.
func (client *Client) Complete(
	ctx context.Context,
	model types.Model,
	prompt string,
) (types.Response, error) {
	return client.Backend.Complete(ctx, model, prompt)
}

// Chat initiates a chat conversation with the provided chat model. Returns a
// Conversation object with which messages can be sent and received.
func (client *Client) Chat(model types.Model) types.Conversation {
	return client.Backend.Chat(model)
}

// GenerateCode sends the provided prompt to the backend and returns a
// Response object. It is a convenience wrapper around client.Complete (for
// text completion models) and client.Chat.Send (for chat models).
func (client *Client) GenerateCode(
	ctx context.Context,
	model types.Model,
	prompt string,
	msgs ...types.Message,
) (res types.Response, err error) {
	if model.Type == types.ModelTypeChat {
		chat := client.Chat(model)
		return chat.Send(ctx, prompt, msgs...)
	}

	return client.Complete(ctx, model, prompt)
}
