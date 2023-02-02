package collectioncallgrpc

import (
	"context"
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Pagination_GRPC(t *testing.T) {
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
	)

	// bug cannot select with flexibility
	testCases := []struct {
		name string
		exec func(collectionxclient.Collector) collectionxclient.Collector
	}{
		{
			name: "Forward 2 Docs To Page 1",
			exec: func(query collectionxclient.Collector) collectionxclient.Collector {
				return query.Page(1).Limit(2).OrderBy("created_at", collectionxclient.Asc)
			},
		}, {
			name: "Forward 2 Docs To Page 3",
			exec: func(query collectionxclient.Collector) collectionxclient.Collector {
				return query.Page(3).Limit(2).OrderBy("created_at", collectionxclient.Asc)
			},
		}, {
			name: "Backward 2 Docs To Page 2",
			exec: func(query collectionxclient.Collector) collectionxclient.Collector {
				return query.Page(2).Limit(2).OrderBy("created_at", collectionxclient.Asc)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			q := tc.exec(query)
			res, err := q.Retrive()
			if err != nil {
				t.Error(err)
			}

			spew.Dump(res)
		})

	}
}
