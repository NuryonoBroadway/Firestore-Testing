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

// List Firestore All Documents ASC
func Test_Firestore_Find_All_Documents_Sorting_ASC_Pagination(t *testing.T) {
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

		// set sorting in here with order by
		query = doc.Collection("cities").OrderBy("name", firestore.Asc)
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

	limit := 2
	firstIndex := 0
	lastIndex := limit
	testCase := []struct {
		name    string
		execute func()
	}{
		{
			name: "Forward 2 Steps",
			execute: func() {

				spew.Dump(cities[firstIndex:lastIndex])
			},
		},
		{
			name: "Forward 2 Steps",
			execute: func() {
				firstIndex = limit + 1
				lastIndex = limit + limit + 1
				spew.Dump(cities[firstIndex:lastIndex])
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.execute()
		})
	}

}
