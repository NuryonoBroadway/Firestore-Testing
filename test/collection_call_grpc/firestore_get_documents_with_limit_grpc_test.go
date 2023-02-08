package collectioncallgrpc

import (
	"context"
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func Test_Get_Documents_With_Limit_GRPC(t *testing.T) {
	cfg, err := collectionxclient.NewClientConfig(
		collectionxclient.WithGrpcAddress("0.0.0.0:9090"),
		collectionxclient.WithProjectRootCollection("development-privypass_collection-core-se"),
		collectionxclient.WithProjectRootDocuments("default"),
		collectionxclient.WithPubSubTopic("test-api"),
		collectionxclient.WithProjectName("cellular-effect-306806"),
	)
	if err != nil {
		t.Error(err)
	}

	var (
		collection_core_client = collectionxclient.NewCollectionClient(cfg)
		ctx, cancel            = context.WithTimeout(context.Background(), 5*time.Second)
	)
	defer cancel()

	conn, err := collection_core_client.OpenConnection(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	var (
		query = conn.Col("development-privypass_collection-core-se").Doc("default").Col("root-collection-test").Doc("default").Col("cities")
	)
	limit := 1
	res, err := query.Limit(limit).Retrive()
	if err != nil {
		t.Error(err)
	}
	spew.Dump(res)
}
