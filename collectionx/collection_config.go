package collectionx

import (
	"context"
	"firebaseapi/helper"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type Config struct {
	GrpcAddress        string `json:"grpc_address"`
	FirebasePath       string `json:"firebase_path"`
	PubSubTopic        string `json:"pubsub_topic"`
	ExternalCollection string `json:"external_collection"`
	ExternalDocument   string `json:"externale_document"`
	ProjectID          string `json:"project_id"`
}

type StandardAPI struct {
	Status  string `json:"status,omitempty"`
	Entity  string `json:"entity,omitempty"`
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`

	// Note :
	//	- you must check a type
	// 	- can be a json or array json
	//	- can be a []map[string]interface{} or map[string]interface{}
	Data  Data   `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type Data struct {
	Type string
	Data interface{}
}

type Error struct {
	General    string              `json:"general"`
	Validation []map[string]string `json:"validation"`
}

type ListValue struct {
	RefID  string                 `json:"ref_id"`
	Object map[string]interface{} `json:"object"`
}

func registryFirestoreClient(cfg Config) *firestore.Client {
	ctx := context.Background()
	credOpt := option.WithCredentialsFile(cfg.FirebasePath)
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

func (s *StandardAPI) MapValue() map[string]interface{} {
	switch s.Data.Type {
	case helper.Collection:
		if v, ok := s.Data.Data.([]ListValue); ok {
			mapped := map[string]interface{}{}
			for _, v := range v {
				mapped[v.RefID] = v.Object
			}

			return mapped
		}
	case helper.Document:
		if v, ok := s.Data.Data.(map[string]interface{}); ok {
			return v
		}
	}

	return nil
}
