package collectionxserver

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

var (
	pongWait = 60 * time.Second
	// ping interval harus di set kurang dari pongWait, untuk mendapatkan 90% dari waktu gunakan 0 * 9 / 10
	// alasanya adalah pingInterval harus lebih rendah dari pingRequecnty, dikarenakan nantinya akan send ping baru sebelum mendapatkan response
	pingInterval = (pongWait * 9) / 10
)

type Subscriber interface {
	Subscribe(ctx context.Context, handler func(context.Context, *Message) error, opts ...Opts) error
}

type consumer struct {
	logger  *logrus.Logger
	cfg     *Options
	client  *pubsub.Client
	source  CollectionCore_SourceDocument
	message chan *pubsub.Message
}

func NewConsumer(cfg *Options, source CollectionCore_SourceDocument, client *pubsub.Client) *consumer {
	return &consumer{
		logger:  logrus.New(),
		cfg:     cfg,
		source:  source,
		client:  client,
		message: make(chan *pubsub.Message),
	}
}

func (c *consumer) Read(ctx context.Context) {
	defer func() {
		c.client.Close()
	}()

	sub := c.client.Subscription(c.cfg.SubscribeID)
	sub.ReceiveSettings.Synchronous = c.cfg.SubscribeAsync
	sub.ReceiveSettings.MaxOutstandingMessages = c.cfg.MaxConcurrent

	err := sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		c.message <- msg
	})

	if err != nil {
		c.logger.Fatal("read subs error: ", err)
	}
}

func (c *consumer) Write(ctx context.Context) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		c.client.Close()
	}()

	for {
		select {
		case msg, ok := <-c.message:
			if ok {
				var data *Payload
				if err := json.Unmarshal(msg.Data, &data); err != nil {
					msg.Nack()
					break
				}

				if err := c.source.Save(ctx, data); err != nil {
					msg.Nack()
					break
				}

				c.logger.Infof("success commited doc: %v", data.Path[len(data.Path)-1].DocumentID)
				msg.Ack()
			}

		case <-ticker.C:
			// keep it stay alive
		}
	}
}
