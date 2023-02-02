package collectioncallgrpc

import (
	"context"
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"firebaseapi/helper"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Get_All_Documents_Collection_Group_GRPC(t *testing.T) {
	cfg, err := collectionxclient.NewClientConfig(
		collectionxclient.WithGrpcAddress("0.0.0.0:9090"),
		collectionxclient.WithProjectRootCollection("development-privypass_collection-core-se"),
		collectionxclient.WithProjectRootDocuments("default"),
		collectionxclient.WithPubSubTopic("pubsub"),
	)
	if err != nil {
		t.Error(err)
	}

	var (
		collection_core_client = collectionxclient.NewCollectionClient(cfg)
		ctx                    = context.Background()
	)

	conn, err := collection_core_client.OpenConnection(ctx)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	var (
		query = conn.ColGroup("cities").Where("country", helper.EqualTo, "USA")
	)

	res, err := query.Retrive()
	if err != nil {
		t.Error(err)
	}
	spew.Dump(res)

}
