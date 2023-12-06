package bedrock

import (
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// Client is the struct that implements libaiac's Backend interface.
type Client struct {
	backend *bedrockruntime.Client
}

// NewClient constructs a new Client object. It receives the standard Options
// object from the AWS SDK for the bedrockruntime service.
func NewClient(opts bedrockruntime.Options) *Client {
	return &Client{
		backend: bedrockruntime.New(opts),
	}
}
