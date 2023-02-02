package collectioncallgrpc

import (
	"context"
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Get_Documents_With_Sort_GRPC(t *testing.T) {
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
		// main_col = collectionx.NewCollectionPayloads(collectionx.WithRootCollection(config.ExternalCollection))
		query = conn.Col("root-collection-test").Doc("default").Col("cities")
		// query = conn.Doc("default").Col("root-collection-test").Doc("default").Col("cities")
	)

	sorts := Sort{
		By:  "name",
		Dir: collectionxclient.Asc,
	}

	res, err := query.OrderBy(sorts.By, sorts.Dir).Retrive()
	if err != nil {
		t.Error(err)
	}
	spew.Dump(res)
}
