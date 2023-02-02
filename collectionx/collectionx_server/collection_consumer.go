package collectionxserver

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type Subscriber interface {
	Subscribe(ctx context.Context, handler func(context.Context, *Message) error, opts ...Opts) error
}

type consumer struct {
	cfg    *ServerConfig
	client *pubsub.Client
	source CollectionCore_SourceDocument
}

func NewConsumer(cfg *ServerConfig, source CollectionCore_SourceDocument, client *pubsub.Client) *consumer {
	return &consumer{
		cfg:    cfg,
		source: source,
		client: client,
	}
}

func (p *consumer) Subscribe(ctx context.Context, handler func(context.Context, *Message) error, opts ...Opts) error {
	defer p.client.Close()

	cfg := defaults()
	for i := 0; i < len(opts); i++ {
		opts[i](cfg)
	}

	sub := p.client.Subscription(cfg.Topic)
	sub.ReceiveSettings.Synchronous = cfg.SubscribeAsync
	sub.ReceiveSettings.MaxOutstandingMessages = cfg.MaxConcurrent

	err := sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		if err := handler(ctx, &Message{
			ID:        msg.ID,
			Attribute: msg.Attributes,
			Data:      msg.Data,
		}); err != nil {
			msg.Nack()
		}
		msg.Ack()
	})

	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	return nil
}

func (c *consumer) Processing(ctx context.Context, msg *Message) error {
	var (
		data *Payload
	)

	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("decode message data %v", err)
	}

	if err := c.source.Save(ctx, data); err != nil {
		return fmt.Errorf("save message to firestore %v", err)
	}

	return nil
}
