package collectionxserver

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

type Subscriber interface {
	Subscribe(ctx context.Context, handler func(context.Context, *Message) error, opts ...Opts) error
}

type consumer struct {
	logger *logrus.Logger
	cfg    *ServerConfig
	client *pubsub.Client
	source CollectionCore_SourceDocument
}

func NewConsumer(cfg *ServerConfig, source CollectionCore_SourceDocument, client *pubsub.Client) *consumer {
	return &consumer{
		logger: logrus.New(),
		cfg:    cfg,
		source: source,
		client: client,
	}
}

func (c *consumer) Subscribe(ctx context.Context, opts ...Opts) error {
	defer c.client.Close()

	cfg := defaults()
	for i := 0; i < len(opts); i++ {
		opts[i](cfg)
	}

	sub := c.client.Subscription(cfg.Topic)
	sub.ReceiveSettings.Synchronous = cfg.SubscribeAsync
	sub.ReceiveSettings.MaxOutstandingMessages = cfg.MaxConcurrent

	err := sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		var data *Payload
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			c.logger.Error(err)
			msg.Nack()
		}

		if err := c.source.Save(ctx, data); err != nil {
			c.logger.Error(err)
			msg.Nack()
		}

		c.logger.Info("success commited")
		msg.Ack()
	})

	if err != nil {
		return fmt.Errorf("sub receive: %v", err)
	}

	return nil
}
