package listsubcollection

import (
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func Test_Firestore_List_Sub_Collection(t *testing.T) {
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
		iter  = query.DocumentRefs(ctx)
	)

	// datas := make([]map[string]interface{}, 0)
	var id []string
	for {
		ref, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				fmt.Printf("docs iterator done\n")
				break
			}
			t.Errorf("error docs iteration err: \n%+v\n", err)
			break
		}
		id = append(id, ref.ID)

	}

	spew.Dump(id)
}
