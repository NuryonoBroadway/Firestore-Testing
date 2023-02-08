package collectioncallgrpc

import (
	"context"
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Find_One_Documents_GRPC(t *testing.T) {
	cfg, err := collectionxclient.NewClientConfig(
		collectionxclient.WithGrpcAddress("0.0.0.0:9090"),
		collectionxclient.WithProjectRootCollection("development-privypass_collection-core-se"),
		collectionxclient.WithProjectRootDocuments("default"),
		collectionxclient.WithPubSubTopic("test-api"),
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
		query = conn.Col("development-privypass_collection-core-se").Doc("default").Col("root-collection-test").Doc("default").Col("cities")
	)

	testCases := []struct {
		name  string
		refId string
	}{
		{
			name:  "Find Japan Tokyo",
			refId: "TOK",
		}, {
			name:  "Find China Bejing",
			refId: "BJ",
		}, {
			name:  "Find USA DC",
			refId: "DC",
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			res, err := query.Doc(tc.refId).Retrive()
			if err != nil {
				t.Error(err)
			}
			spew.Dump(res)
		})

	}
}
