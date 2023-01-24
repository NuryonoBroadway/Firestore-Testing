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

// List Firestore All Documents DESC
func Test_Firestore_Find_All_Documents_With_Limit(t *testing.T) {
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
		query = doc.Collection("cities").Limit(2)
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

		var data helper.City
		if err := doc.DataTo(&data); err != nil {
			t.Errorf("data to err: \n%+v\n", err)
			break
		}

		cities = append(cities, data)
	}

	spew.Dump(cities)

}