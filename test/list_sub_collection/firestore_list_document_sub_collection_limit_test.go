package listsubcollection

import (
	"firebaseapi/helper"
	"testing"

	"cloud.google.com/go/firestore"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// List Firestore by ref id With Limit
func Test_Firestore_List_Ref_ID_With_Limit(t *testing.T) {
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

	testCases := []struct {
		name  string
		limit int
	}{
		{
			name:  "With 2 Limit",
			limit: 2,
		}, {
			name:  "With 2 Limit",
			limit: 3,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			iter := query.DocumentRefs(ctx)
			var collRef []helper.City

			reach := 0
			for {
				if reach == tc.limit {
					break
				}

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

				var city helper.City
				if err := queryIterator.DataTo(&city); err != nil {
					t.Error(err)
					break
				}

				collRef = append(collRef, city)

				reach++
			}

			spew.Dump(collRef)
		})
	}

}
