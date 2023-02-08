package collectioncallgrpc

import (
	"context"
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Save_Documents_Update_GRPC(t *testing.T) {
	cfg, err := collectionxclient.NewClientConfig(
		collectionxclient.WithGrpcAddress("0.0.0.0:9090"),
		collectionxclient.WithProjectRootCollection("development-privypass_collection-core-se"),
		collectionxclient.WithProjectRootDocuments("default"),
		collectionxclient.WithProjectName("cellular-effect-306806"),
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
		log.Fatal(err)
	}
	defer conn.Close()
	var (
		query = conn.Col("development-privypass_collection-core-se").Doc("default").Col("root-collection-test").Doc("default").Col("cities").Doc("SF")
		ready = query.Set([]collectionxclient.Row{
			{
				Path:  "capital",
				Value: true,
			}, {
				Path:  "country",
				Value: "USA",
			}, {
				Path:  "name",
				Value: "San Fransisco",
			},
		}, true)
	)

	res, err := ready.Save()
	if err != nil {
		t.Error(err)
	}
	spew.Dump(res)

}
