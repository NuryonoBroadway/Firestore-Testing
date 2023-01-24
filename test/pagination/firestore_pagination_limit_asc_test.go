package pagination

import (
	"testing"

	"cloud.google.com/go/firestore"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/option"
)

// List Firestore Pagination With Limit Cursor With ASC
func Test_Firestore_Pagination_With_Limit_ASC(t *testing.T) {
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

	// construct pagination with flexible steps

	// StartAfter returns a new Query that specifies that results should start just after
	// the document with the given field values. See Query.StartAt for more information.

	// StartAt returns a new Query that specifies that results should start at the document
	// with the given field values. StartAt may be called with a single DocumentSnapshot,
	// representing an existing document within the query. The document must be a direct
	// child of the location being queried (not a parent document, or document in a different
	// collection, or a grandchild document, for example).

	// EndAt returns a new Query that specifies that results should end at the
	// document with the given field values. See Query.StartAt for more information.

	// EndBefore returns a new Query that specifies that results should end just before
	// the document with the given field values. See Query.StartAt for more information.

	// pagination test with cursor goes in here
	cities := []*firestore.DocumentSnapshot{}
	useCase := []struct {
		name      string
		limit     int
		behaviour string
		excute    func(int) error
	}{
		{
			name:      "Forward with 2 steps",
			limit:     2,
			behaviour: "forward",
			excute: func(limit int) error {
				iter := query.OrderBy("population", firestore.Asc).Limit(limit).Documents(ctx)
				// datas := make([]map[string]interface{}, 0)
				docs, err := iter.GetAll()
				if err != nil {
					return err
				}

				for _, v := range docs {
					spew.Dump(v.Data())
				}

				cities = docs
				return nil
			},
		},
		{
			name:      "Backward with 2 steps",
			limit:     2,
			behaviour: "back",
			excute: func(limit int) error {
				lastCity := cities[0]
				iter := query.OrderBy("population", firestore.Asc).StartAt(lastCity.Data()["population"]).Limit(limit).Documents(ctx)
				docs, err := iter.GetAll()
				if err != nil {
					return err
				}

				for _, v := range docs {
					spew.Dump(v.Data())
				}

				cities = docs

				return nil
			},
		},
		{
			name:      "Forward with 2 steps",
			limit:     2,
			behaviour: "forward",
			excute: func(limit int) error {
				lastCity := cities[len(cities)-1]
				iter := query.OrderBy("population", firestore.Asc).StartAfter(lastCity.Data()["population"]).Limit(limit).Documents(ctx)
				docs, err := iter.GetAll()
				if err != nil {
					return err
				}

				for _, v := range docs {
					spew.Dump(v.Data())
				}
				cities = docs

				return nil
			},
		},
		{
			name:      "Forward with 2 steps",
			limit:     2,
			behaviour: "forward",
			excute: func(limit int) error {
				lastCity := cities[len(cities)-1]
				iter := query.OrderBy("population", firestore.Asc).StartAfter(lastCity.Data()["population"]).Limit(limit).Documents(ctx)
				docs, err := iter.GetAll()
				if err != nil {
					return err
				}

				for _, v := range docs {
					spew.Dump(v.Data())
				}
				cities = docs

				return nil
			},
		},
	}

	for i := range useCase {
		tc := useCase[i]

		t.Run(tc.name, func(t *testing.T) {
			if err := tc.excute(tc.limit); err != nil {
				t.Error(err)
				return
			}
		})
	}
}
