package collectionxclient

import (
	"github.com/go-playground/validator/v10"
)

type (
	ClientConfig struct {
		GrpcAddress           string `json:"grpc_address" validate:"required"`
		PubSubTopic           string `json:"pubsub_topic"`
		ProjectRootCollection string `json:"project_root_collection" validate:"required"`
		ProjectRootDocument   string `json:"project_root_document"`
	}

	Option func(p *ClientConfig)
)

var validate = validator.New()

func NewClientConfig(opts ...Option) (*ClientConfig, error) {
	p := ClientConfig{}
	for i := 0; i < len(opts); i++ {
		opts[i](&p)
	}

	if err := validate.Struct(&p); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			return nil, e
		}
	}

	return &p, nil
}

func WithGrpcAddress(address string) Option {
	return func(c *ClientConfig) {
		c.GrpcAddress = address
	}
}

func WithPubSubTopic(topic string) Option {
	return func(c *ClientConfig) {
		c.PubSubTopic = topic
	}
}

func WithProjectRootCollection(col string) Option {
	return func(c *ClientConfig) {
		c.ProjectRootCollection = col
	}
}

func WithProjectRootDocuments(doc string) Option {
	return func(c *ClientConfig) {
		c.ProjectRootDocument = doc
	}
}
