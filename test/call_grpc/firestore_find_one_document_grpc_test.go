package callgrpc

import (
	"firebaseapi/collectionx"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Find_One_Documents_GRPC(t *testing.T) {
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
			spew.Dump(res.MapValue())
		})

	}
}
