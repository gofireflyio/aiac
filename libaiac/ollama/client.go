package ollama

import ol "github.com/jmorganca/ollama/api"

// Client is the struct that implements libaiac's Backend interface.
type Client struct {
	backend *ol.Client
}

func NewClient() *Client {
	return &Client{
		backend: newBackend(),
	}
}

func newBackend() *ol.Client {
	client, _ := ol.ClientFromEnvironment()
	return client
}
