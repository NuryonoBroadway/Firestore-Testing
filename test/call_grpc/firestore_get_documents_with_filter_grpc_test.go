package callgrpc

import (
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"firebaseapi/helper"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Get_Documents_With_Filter_GRPC(t *testing.T) {
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
	)

	conn, err := collection_core_client.OpenConnection()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	var (
		// main_col = collectionx.NewCollectionPayloads(collectionx.WithRootCollection(config.ExternalCollection))
		query = conn.Col("root-collection-test").Doc("default").Col("cities")
		// query = conn.Doc("default").Col("root-collection-test").Doc("default").Col("cities")
	)

	filters := []collectionxclient.Filter{
		{
			By:  "country",
			Op:  helper.EqualTo,
			Val: "USA",
		},
		{
			By:  "capital",
			Op:  helper.EqualTo,
			Val: false,
		},
	}

	res, err := query.Where(filters...).Retrive()
	if err != nil {
		t.Error(err)
	}
	spew.Dump(res.MapValue())
}
