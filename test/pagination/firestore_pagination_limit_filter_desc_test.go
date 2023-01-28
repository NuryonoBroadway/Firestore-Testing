package pagination

import (
	"firebaseapi/helper"
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/option"
)

// List Firestore Pagination With Limit Cursor With ASC
func Test_Firestore_Pagination_With_Filter_Limit_DESC(t *testing.T) {
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

		collect = doc.Collection("cities").Query
	)

	// pagination test with cursor goes in here
	type filter struct {
		Operator helper.Operator
		Value    interface{}
	}

	// needs index
	cities := []*firestore.DocumentSnapshot{}
	useCase := []struct {
		name      string
		limit     int
		behaviour string
		filter    map[string]filter
		excute    func(int, map[string]filter) error
	}{
		{
			name:      "Forward with 2 steps With Capital True",
			limit:     2,
			behaviour: "forward",
			filter: map[string]filter{
				"capital": {
					Operator: helper.EqualTo,
					Value:    true,
				},
			},
			excute: func(limit int, filter map[string]filter) error {
				// iteration is empty or can iterate the multiple filters
				// index by is important to use in case using > < condition
				for k, v := range filter {
					collect = collect.Where(k, v.Operator.ToString(), v.Value)
				}

				iter := collect.OrderBy("population", firestore.Desc).Limit(limit).Documents(ctx)
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
			name:      "Forward with 2 steps With Capital True",
			limit:     2,
			behaviour: "forward",
			filter: map[string]filter{
				"capital": {
					Operator: helper.EqualTo,
					Value:    true,
				},
			},
			excute: func(limit int, filter map[string]filter) error {
				lastCity := cities[len(cities)-1]

				// iteration is empty or can iterate the multiple filters
				// index by is important to use in case using > < condition
				for k, v := range filter {
					collect = collect.Where(k, v.Operator.ToString(), v.Value)
				}

				iter := collect.OrderBy("population", firestore.Desc).StartAfter(lastCity.Data()["population"]).Limit(limit).Documents(ctx)
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
			name:      "Backward with 2 steps With Capital True",
			limit:     2,
			behaviour: "backward",
			filter: map[string]filter{
				"capital": {
					Operator: helper.EqualTo,
					Value:    true,
				},
			},
			excute: func(limit int, filter map[string]filter) error {
				lastCity := cities[0]
				fmt.Println("last city ", lastCity.Data())
				collect := doc.Collection("cities").Query

				// iteration is empty or can iterate the multiple filters
				// index by is important to use in case using > < condition
				for k, v := range filter {
					collect = collect.Where(k, v.Operator.ToString(), v.Value)
				}

				iter := collect.OrderBy("population", firestore.Asc).StartAfter(lastCity.Data()["population"]).Limit(limit).Documents(ctx)
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
	}

	for i := range useCase {
		tc := useCase[i]

		t.Run(tc.name, func(t *testing.T) {
			if err := tc.excute(tc.limit, tc.filter); err != nil {
				t.Error(err)
				return
			}
		})
	}
}
