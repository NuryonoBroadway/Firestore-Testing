package collectioncallgrpc

import (
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"firebaseapi/helper"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Pagination_With_Filter_GRPC(t *testing.T) {
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

	testCases := []struct {
		name   string
		filter []Filter
		exec   func(query collectionxclient.Collector, filter []Filter) collectionxclient.Collector
	}{
		{
			name: "Forward 2 Docs",
			filter: []Filter{
				{
					By:  "capital",
					Op:  helper.EqualTo,
					Val: false,
				},
			},
			exec: func(query collectionxclient.Collector, filter []Filter) collectionxclient.Collector {
				for _, v := range filter {
					query = query.Where(v.By, v.Op, v.Val)
				}
				return query.Page(1).Limit(2).OrderBy("created_at", collectionxclient.Asc)
			},
		},
		{
			name: "Backward 2 Docs",
			filter: []Filter{
				{
					By:  "capital",
					Op:  helper.EqualTo,
					Val: false,
				},
			},
			exec: func(query collectionxclient.Collector, filter []Filter) collectionxclient.Collector {
				for _, v := range filter {
					query = query.Where(v.By, v.Op, v.Val)
				}
				return query.Page(1).Limit(2).OrderBy("created_at", collectionxclient.Asc)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			q := tc.exec(query, tc.filter)
			res, err := q.Retrive()
			if err != nil {
				t.Error(err)
			}
			spew.Dump(res)
		})

	}
}
