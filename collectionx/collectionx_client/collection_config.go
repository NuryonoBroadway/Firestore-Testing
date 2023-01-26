package collectionxclient

import (
	"github.com/go-playground/validator/v10"
)

type configClient struct {
	GrpcAddress           string `json:"grpc_address" validate:"required"`
	PubSubTopic           string `json:"pubsub_topic"`
	ProjectRootCollection string `json:"project_root_collection" validate:"required"`
	ProjectRootDocument   string `json:"project_root_document"`
}

type clientCollector func(p *configClient)

func NewClientConfig(opts ...clientCollector) (*configClient, error) {
	validate := validator.New()

	p := configClient{}
	for _, v := range opts {
		v(&p)
	}

	if err := validate.Struct(&p); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			return nil, e
		}
	}

	return &p, nil
}

func WithGrpcAddress(address string) clientCollector {
	return func(c *configClient) {
		c.GrpcAddress = address
	}
}

func WithPubSubTopic(pubsub string) clientCollector {
	return func(c *configClient) {
		c.PubSubTopic = pubsub
	}
}

func WithProjectRootCollection(collection string) clientCollector {
	return func(c *configClient) {
		c.ProjectRootCollection = collection
	}
}

func WithProjectRootDocuments(documents string) clientCollector {
	return func(c *configClient) {
		c.ProjectRootDocument = documents
	}
}
