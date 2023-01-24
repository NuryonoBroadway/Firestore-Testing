package betweensnapshots

import (
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Listen Snapshots
func Test_Firestore_Changes_Between_Snapshots(t *testing.T) {
	client, err := firestore.NewClient(ctx, project_id, option.WithCredentialsFile(credential_file))
	if err != nil {
		t.Errorf("failed open firestore client err: \n%+v\n", err)
	}

	defer client.Close()

	var (
		main_col = client.Collection(collection_id)
		main_doc = main_col.Doc("default")
	)

	var (
		col = main_doc.Collection("root-collection-test")
		doc = col.Doc("default")

		// set limit in here
		query = doc.Collection("cities")
		iter  = query.Snapshots(ctx)
	)

	defer iter.Stop()

	for {
		snap, err := iter.Next()
		if err != nil {
			if status.Code(err) == codes.DeadlineExceeded {
				return
			}
			t.Error(t)
			return
		}

		if snap != nil {
			for _, change := range snap.Changes {
				switch change.Kind {
				case firestore.DocumentAdded:
					fmt.Printf("New: %v", change.Doc.Data())
				case firestore.DocumentModified:
					fmt.Printf("Modified: %v", change.Doc.Data())
				case firestore.DocumentRemoved:
					fmt.Printf("Modified: %v", change.Doc.Data())
				}
			}
		}

	}
}
