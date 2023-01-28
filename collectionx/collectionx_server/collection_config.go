package collectionxserver

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type ServerConfig struct {
	ProjectID       string `json:"project_id"`
	CredentialsFile string `json:"credentials_file"`
	RootCollection  string `json:"root_collection"`
	RootDocument    string `json:"root_document"`
	ProjectServiceConfig
}

type ProjectServiceConfig struct {
	Environment string `json:"environment"`
	ServiceName string `json:"service_name"`
	ProjectName string `json:"project_name"`
}

func RegistryFirestoreClient(cfg *ServerConfig) *firestore.Client {
	ctx := context.Background()
	credOpt := option.WithCredentialsFile(cfg.CredentialsFile)
	conf := &firebase.Config{ProjectID: cfg.ProjectID}
	app, err := firebase.NewApp(ctx, conf, credOpt)
	if err != nil {
		logger.Fatal(fmt.Sprintf("google firestore error:%v", err))
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		logger.Fatal(fmt.Sprintf("google firestore error:%v", err))
	}

	return client
}
