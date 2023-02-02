package collectionxserver

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
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
