package collectioncallgrpc

import (
	"context"
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"log"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func Test_Save_Documents_Create_GRPC(t *testing.T) {
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
		query = conn.Col("development-privypass_collection-core-se").Doc("default").Col("root-collection-test").Doc("default").Col("cities").Doc("JKT")
		ready = query.Set([]collectionxclient.Row{
			{
				Path:  "country",
				Value: "Indonesia",
			}, {
				Path:  "capital",
				Value: true,
			}, {
				Path:  "name",
				Value: "Jakarta",
			}, {
				Path:  "Population",
				Value: 89000000000,
			}, {
				Path:  "created_at",
				Value: time.Now(),
			}, {
				Path:  "updated_at",
				Value: time.Now(),
			},
		}, true)
	)

	res, err := ready.Save()
	if err != nil {
		t.Error(err)
	}
	spew.Dump(res)

}
