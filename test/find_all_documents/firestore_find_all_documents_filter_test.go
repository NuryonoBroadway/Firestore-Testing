package findalldocuments

import (
	"firebaseapi/helper"
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Find All Documents With Filtering
func Test_Firestore_FindAllDocumentsFilter(t *testing.T) {
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

		query = doc.Collection("cities").Where("capital", "==", true)
		iter  = query.Documents(ctx)
	)

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
