package collectionxserver

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/grpc/keepalive"
)

type ServerConfig struct {
	FirestoreProjectID       string `json:"firestore_project_id"`
	FirestoreCredentialsFile string `json:"firestore_credentials_file"`
	PubSubProjectID          string `json:"pubsub_project_id"`
	PubSubCredentialsFile    string `json:"pubsub_credentials_file"`
	RootCollection           string `json:"root_collection"`
	RootDocument             string `json:"root_document"`
}

func (cfg *ServerConfig) RegistryFirestoreClient(ctx context.Context) *firestore.Client {
	credOpt := option.WithCredentialsFile(cfg.FirestoreCredentialsFile)
	conf := &firebase.Config{ProjectID: cfg.FirestoreProjectID}
	app, err := firebase.NewApp(ctx, conf, credOpt)
	if err != nil {
		logger.Fatal(fmt.Sprintf("google firestore error:%v", err)) // removed
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		logger.Fatal(fmt.Sprintf("google firestore error:%v", err))
	}

	return client
}

func (cfg *ServerConfig) RegistryPubSubConsumer(ctx context.Context) *pubsub.Client {
	credOpt := option.WithCredentialsFile(cfg.PubSubCredentialsFile)
	client, err := pubsub.NewClient(ctx, cfg.PubSubProjectID, credOpt)
	if err != nil {
		logger.Fatal(fmt.Sprintf("google pusbsub conusmer error:%v", err))
	}

	return client
}

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}
