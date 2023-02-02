package collectiongroup

import (
	"firebaseapi/helper"
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Find All Documents With Sorting ASC
func Test_Firestore_Collection_Group_GRPC(t *testing.T) {
	client, err := firestore.NewClient(ctx, project_id, option.WithCredentialsFile(credential_file))
	if err != nil {
		t.Errorf("failed open firestore client err: \n%+v\n", err)
	}

	defer client.Close()

	// collection group dont need to iterate over the collection
	// just choose available collection
	iter := client.CollectionGroup("cities").OrderBy("created_at", firestore.Asc).Documents(ctx)
	defer iter.Stop()

	cities := []helper.City{}
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				fmt.Printf("docs iterator done\n")
				break
			}
			t.Errorf("error docs iteration err: \n%+v\n", err)
			break
		}

		var city helper.City
		if err := doc.DataTo(&city); err != nil {
			t.Errorf("data to err: \n%+v\n", err)
			break
		}

		cities = append(cities, city)
	}

	spew.Dump(cities)
}
