package libaiac

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gofireflyio/aiac/v5/libaiac/bedrock"
	"github.com/gofireflyio/aiac/v5/libaiac/ollama"
	"github.com/gofireflyio/aiac/v5/libaiac/openai"
	"github.com/gofireflyio/aiac/v5/libaiac/types"
)

// Version contains aiac's version string
var Version = "development"

// Aiac provides the main interface for using libaiac.
type Aiac struct {
	// Conf holds the configuration for aiac.
	Conf Config

	// Backends is a map from backend names to backend implementations.
	Backends map[string]types.Backend
}

// New constructs a new Aiac object with the path to a configuration file. If
// a configuration file is not provided, the default path will be checked based
// on the XDG specification. On Unix-like operating systems, this will be
// ~/.config/aiac/aiac.toml.
func New(configPath ...string) (*Aiac, error) {
	path := ""
	if len(configPath) > 0 {
		path = configPath[0]
	}

	conf, err := LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("failed loading configuration: %w", err)
	}

	return &Aiac{Conf: conf}, nil
}

// NewFromConf is the same as New, but receives a populated configuration object
// rather than a file path.
func NewFromConf(conf Config) *Aiac {
	return &Aiac{Conf: conf}
}

// ListModels returns a list of all the models supported by the selected
// backend, identified by its name. If backendName is an empty string, the
// default backend defined in the configuration file will be used, if any.
func (aiac *Aiac) ListModels(ctx context.Context, backendName string) (
	models []string,
	err error,
) {
	backend, _, err := aiac.loadBackend(ctx, backendName)
	if err != nil {
		return models, fmt.Errorf("failed loading backend: %w", err)
	}

	return backend.ListModels(ctx)
}

// Chat initiates a chat conversation with the provided chat model of the
// selected backend. Returns a Conversation object with which messages can be
// sent and received. If backendName is an empty string, the default backend
// defined in the configuration will be used, if any. If model is an empty
// string, the default model defined in the backend configuration will be used,
// if any. Users can also supply zero or more "previous messages" that may have
// been exchanged in the past. This practically allows "loading" previous
// conversations and continuing them.
func (aiac *Aiac) Chat(
	ctx context.Context,
	backendName string,
	model string,
	msgs ...types.Message,
) (chat types.Conversation, err error) {
	backend, defaultModel, err := aiac.loadBackend(ctx, backendName)
	if err != nil {
		return chat, fmt.Errorf("failed loading backend: %w", err)
	}

	if model == "" {
		if defaultModel == "" {
			return nil, types.ErrNoDefaultModel
		}
		model = defaultModel
	}

	return backend.Chat(model, msgs...), nil
}

func (aiac *Aiac) loadBackend(ctx context.Context, name string) (
	backend types.Backend,
	defaultModel string,
	err error,
) {
	if name == "" {
		if aiac.Conf.DefaultBackend == "" {
			return nil, defaultModel, types.ErrNoDefaultBackend
		}
		name = aiac.Conf.DefaultBackend
	}

	// Check if we've already loaded it before
	if backend, ok := aiac.Backends[name]; ok {
		return backend, defaultModel, nil
	}

	// We haven't, check if it's in the configuration
	backendConf, ok := aiac.Conf.Backends[name]
	if !ok {
		return backend, defaultModel, types.ErrNoSuchBackend
	}

	switch backendConf.Type {
	case BackendBedrock:
		if backendConf.AWSProfile == "" {
			backendConf.AWSProfile = bedrock.DefaultAWSProfile
		}

		if backendConf.AWSRegion == "" {
			backendConf.AWSRegion = bedrock.DefaultAWSRegion
		}

		cfg, err := config.LoadDefaultConfig(
			ctx,
			config.WithSharedConfigProfile(backendConf.AWSProfile),
		)
		if err != nil {
			return nil, defaultModel, err
		}

		cfg.Region = backendConf.AWSRegion

		backend = bedrock.New(cfg)
	case BackendOllama:
		backend = ollama.New(&ollama.Options{
			URL:          backendConf.URL,
			ExtraHeaders: backendConf.ExtraHeaders,
		})
	default:
		// default to openai
		backend, err = openai.New(&openai.Options{
			ApiKey:       backendConf.APIKey,
			URL:          backendConf.URL,
			APIVersion:   backendConf.APIVersion,
			ExtraHeaders: backendConf.ExtraHeaders,
		})
		if err != nil {
			return nil, defaultModel, err
		}
	}

	return backend, backendConf.DefaultModel, nil
}
