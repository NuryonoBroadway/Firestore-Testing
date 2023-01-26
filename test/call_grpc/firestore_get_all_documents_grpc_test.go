package callgrpc

import (
	"firebaseapi/collectionx"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Get_All_Documents_GRPC(t *testing.T) {
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

	res, err := query.Retrive()
	if err != nil {
		t.Error(err)
	}
	spew.Dump(res.MapValue())

}
