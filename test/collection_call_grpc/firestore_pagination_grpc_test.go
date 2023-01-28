package collectioncallgrpc

import (
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
	)

	conn, err := collection_core_client.OpenConnection()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	var (
		// main_col = collectionx.NewCollectionPayloads(collectionx.WithRootCollection(config.ExternalCollection))
		query = conn.Col("root-collection-test").Doc("default").Col("cities")
	)

	docs := collectionxclient.MetaData{}

	testCases := []struct {
		name string
		exec func(collectionxclient.Collector, collectionxclient.MetaData) collectionxclient.Collector
	}{
		{
			name: "Forward 2 Docs",
			exec: func(query collectionxclient.Collector, docs collectionxclient.MetaData) collectionxclient.Collector {
				return query.Pagination(1, docs).Limit(2).OrderBy("population", collectionxclient.Asc)
			},
		},
		{
			name: "Forward 2 Docs",
			exec: func(query collectionxclient.Collector, docs collectionxclient.MetaData) collectionxclient.Collector {
				return query.Pagination(2, docs).Limit(2).OrderBy("population", collectionxclient.Asc)
			},
		},
		{
			name: "Forward 2 Docs",
			exec: func(query collectionxclient.Collector, docs collectionxclient.MetaData) collectionxclient.Collector {
				return query.Pagination(3, docs).Limit(2).OrderBy("population", collectionxclient.Asc)
			},
		},
		{
			name: "Backward 2 Docs",
			exec: func(query collectionxclient.Collector, docs collectionxclient.MetaData) collectionxclient.Collector {
				return query.Pagination(2, docs).Limit(2).OrderBy("population", collectionxclient.Asc)
			},
		},
		{
			name: "Backward 2 Docs",
			exec: func(query collectionxclient.Collector, docs collectionxclient.MetaData) collectionxclient.Collector {
				return query.Pagination(1, docs).Limit(2).OrderBy("population", collectionxclient.Asc)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			q := tc.exec(query, docs)
			res, err := q.Retrive()
			if err != nil {
				t.Error(err)
			}

			response := res.MapValue()
			docs = collectionxclient.MetadataCreator(res.Meta.Page, response)
			spew.Dump(response)
		})

	}
}
