package createfirestore

import (
	"firebaseapi/helper"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func Test_Firestore_SeedData(t *testing.T) {
	var (
		cities = []struct {
			id string
			c  helper.City
		}{
			{id: "SF", c: helper.City{Name: "San Francisco", State: "CA", Country: "USA", Capital: false, Population: 860000}},
			{id: "LA", c: helper.City{Name: "Los Angeles", State: "CA", Country: "USA", Capital: false, Population: 3900000}},
			{id: "DC", c: helper.City{Name: "Washington D.C.", Country: "USA", Capital: true, Population: 680000}},
			{id: "TOK", c: helper.City{Name: "Tokyo", Country: "Japan", Capital: true, Population: 9000000}},
			{id: "BJ", c: helper.City{Name: "Beijing", Country: "China", Capital: true, Population: 21500000}},
		}
	)

	client, err := firestore.NewClient(ctx, project_id, option.WithCredentialsFile(credential_file))
	if err != nil {
		t.Errorf("failed open firestore client err: \n%+v\n", err)
		return
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

	for _, c := range cities {
		if _, err := doc.Collection("cities").Doc(c.id).Set(ctx, c.c); err != nil {
			t.Error(err)
			return
		}
	}
}
