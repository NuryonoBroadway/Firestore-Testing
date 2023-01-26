package listsubcollection

import (
	"firebaseapi/helper"
	"testing"

	"cloud.google.com/go/firestore"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// List Firestore by ref id With Filters
func Test_Firestore_List_Ref_ID_With_Filter(t *testing.T) {
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

		query = doc.Collection("cities")
	)

	type filter struct {
		Operator string
		Value    interface{}
	}

	testCases := []struct {
		name   string
		filter map[string]filter
	}{
		{
			name: "With 1 Filters",
			filter: map[string]filter{
				"capital": filter{
					Operator: helper.EqualTo,
					Value:    true,
				},
			},
		}, {
			name: "With 2 Filters",
			filter: map[string]filter{
				"country": filter{
					Operator: helper.EqualTo,
					Value:    "USA",
				},
				"population": filter{
					Operator: helper.GreaterThanEqual,
					Value:    10000,
				},
			},
		}, {
			name: "With 1 Filters But With Condition In",
			filter: map[string]filter{
				"country": filter{
					Operator: helper.In,
					Value:    []string{"China", "Japan"},
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			iter := query.DocumentRefs(ctx)
			var collRef []helper.City
			for {
				ref, err := iter.Next()
				if err != nil {
					if err == iterator.Done {
						break
					}
					t.Errorf("error docs iteration err: \n%+v\n", err)
					break
				}

				queryIterator, err := query.Doc(ref.ID).Get(ctx)
				if err != nil {
					if err == iterator.Done {
						break
					}
					t.Errorf("error docs iteration err: \n%+v\n", err)
					break
				}

				pass := true
				for i, v := range tc.filter {
					switch v.Operator {
					case helper.EqualTo:
						pass = queryIterator.Data()[i] == v.Value

					case helper.NotEqualTo:
						pass = queryIterator.Data()[i] != v.Value

					case helper.GreaterThan:
						pass = int(queryIterator.Data()[i].(int64)) > v.Value.(int)

					case helper.GreaterThanEqual:
						pass = int(queryIterator.Data()[i].(int64)) >= v.Value.(int)

					case helper.LessThan:
						pass = int(queryIterator.Data()[i].(int64)) < v.Value.(int)

					case helper.LessThanEqual:
						pass = int(queryIterator.Data()[i].(int64)) <= v.Value.(int)

					case helper.In:
						pass = helper.SliceCheckCondition(v.Value, queryIterator.Data()[i])

					case helper.NotIn:
						pass = !helper.SliceCheckCondition(v.Value, queryIterator.Data()[i])
					}

					if !pass {
						break
					}
				}

				if pass {
					var city helper.City
					if err := queryIterator.DataTo(&city); err != nil {
						t.Error(err)
						break
					}

					collRef = append(collRef, city)
				}

			}

			spew.Dump(collRef)
		})
	}

}
