package collectioncallgrpc

import (
	"context"
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"io"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func Test_Collection_Snapshots_GRPC(t *testing.T) {
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

	snap, err := query.Snapshots()
	if err != nil {
		t.Error(err)
	}

	defer snap.Close()
	for {
		res, err := snap.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}

			t.Error(err)
		}

		switch res.Kind {
		case collectionxclient.DOCUMENT_KIND_ADDED.ToString():
			logrus.Info("document added")
			spew.Dump(res)
		case collectionxclient.DOCUMENT_KIND_REMOVED.ToString():
			logrus.Info("document removed")
			spew.Dump(res)
		case collectionxclient.DOCUMENT_KIND_MODIFIED.ToString():
			logrus.Info("document modified")
			spew.Dump(res)
		}

	}
}
