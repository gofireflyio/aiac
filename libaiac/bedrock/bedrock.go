package bedrock

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// Bedrock is the struct that implements libaiac's Backend interface.
type Bedrock struct {
	runtime *bedrockruntime.Client
	service *bedrock.Client
}

const (
	// DefaultAWSRegion is the default AWS region to use if the backend does not
	// specify one
	DefaultAWSRegion = "us-east-1"

	// DefaultAWSProfile is the default AWS profile to use if the backend does
	// not specify one
	DefaultAWSProfile = "default"
)

// New constructs a new Bedrock object. It receives a standard aws.Config
// object.
func New(cfg aws.Config) *Bedrock {
	return &Bedrock{
		runtime: bedrockruntime.NewFromConfig(cfg),
		service: bedrock.NewFromConfig(cfg),
	}
}
