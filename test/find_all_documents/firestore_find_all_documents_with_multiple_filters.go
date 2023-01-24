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

func Test_Firestore_Find_All_Documents_Multiple_Filter(t *testing.T) {
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
	)

	type filter struct {
		Operator string
		Value    interface{}
	}

	// use a test case to test many possible response
	testCases := []struct {
		name   string
		filter map[string]filter
		excute func(map[string]filter) *firestore.DocumentIterator
	}{
		{
			name: "Find USA and Capital True",
			filter: map[string]filter{
				"country": {
					Operator: helper.EqualTo,
					Value:    "USA",
				},
				"capital": {
					Operator: helper.EqualTo,
					Value:    false,
				},
			},
			excute: func(filter map[string]filter) *firestore.DocumentIterator {
				// cant use > < failed precondition

				collect := doc.Collection("cities").Query

				// iteration is empty or can iterate the multiple filters
				for k, v := range filter {
					collect = collect.Where(k, v.Operator, v.Value)
				}

				return collect.Documents(ctx)
			},
		},
		{
			name: "Find Capital is false",
			filter: map[string]filter{
				"capital": {
					Operator: helper.EqualTo,
					Value:    false,
				},
			},
			excute: func(filter map[string]filter) *firestore.DocumentIterator {
				// cant use > < failed precondition

				collect := doc.Collection("cities").Query

				// iteration is empty or can iterate the multiple filters
				for k, v := range filter {
					collect = collect.Where(k, v.Operator, v.Value)
				}

				return collect.Documents(ctx)
			},
		},
		{
			name: "Find Capital is false And Population over 9000000",
			filter: map[string]filter{
				"capital": {
					Operator: helper.EqualTo,
					Value:    true,
				},
				"population": {
					Operator: helper.GreaterThanEqual,
					Value:    "9000000",
				},
			},
			excute: func(filter map[string]filter) *firestore.DocumentIterator {
				// cant use > < failed precondition

				collect := doc.Collection("cities").Query

				// iteration is empty or can iterate the multiple filters
				// index by is important to use in case using > < condition
				for k, v := range filter {
					collect = collect.Where(k, v.Operator, v.Value)
				}

				return collect.Documents(ctx)
			},
		},
		{
			name: "Find If Any Object In",
			filter: map[string]filter{
				"country": {
					Operator: helper.In,
					Value:    []string{"USA", "Japan"}, // perform an string array
				},
			},
			excute: func(filter map[string]filter) *firestore.DocumentIterator {
				// cant use > < failed precondition

				collect := doc.Collection("cities").Query

				// iteration is empty or can iterate the multiple filters
				// index by is important to use in case using > < condition
				for k, v := range filter {
					collect = collect.Where(k, v.Operator, v.Value)
				}

				return collect.Documents(ctx)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			iter := tc.excute(tc.filter)
			defer iter.Stop()

			// datas := make([]map[string]interface{}, 0)
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
		})
	}

}
