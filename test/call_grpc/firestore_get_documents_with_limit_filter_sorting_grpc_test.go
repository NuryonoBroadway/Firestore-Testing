package callgrpc

import (
	"firebaseapi/collectionx"
	"firebaseapi/helper"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Get_Documents_With_Limit_Filter_Sort_GRPC(t *testing.T) {
	config := collectionx.Config{
		GrpcAddress:        "0.0.0.0:9090",
		FirebasePath:       "/home/slvr/FirebaseTestingApi/config/privyfellowship-6a4e1-firebase-adminsdk-g0c2d-66e133ed37.json",
		PubSubTopic:        "",
		ExternalCollection: "development-privypass_collection-core-se",
		ExternalDocument:   "default",
		ProjectID:          "privyfellowship-6a4e1",
	}

	var (
		collection_core_client = collectionx.NewClient(&config)
	)

	conn, err := collection_core_client.OpenConnection()
	if err != nil {
		t.Error(err)
	}

	defer conn.Close()

	var (
		// main_col = collectionx.NewCollectionPayloads(collectionx.WithRootCollection(config.ExternalCollection))
		query = conn.Doc("default").Col("root-collection-test").Doc("default").Col("cities")
		// query = conn.Doc("default").Col("root-collection-test").Doc("default").Col("cities")
	)

	testCases := []struct {
		name    string
		filters []collectionx.Filter
		sorts   []collectionx.Sort
		limit   int
	}{
		{
			name: "With filter",
			filters: []collectionx.Filter{
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
			},
		}, {
			name: "With Sort",
			sorts: []collectionx.Sort{
				{
					By:  "name",
					Dir: helper.ASC,
				},
			},
		},
		{
			name:  "With Limit",
			limit: 1,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			res, err := query.Where(tc.filters...).Order(tc.sorts...).Max(tc.limit).Retrive()
			if err != nil {
				t.Error(err)
			}
			spew.Dump(res.MapValue())
		})

	}
}
