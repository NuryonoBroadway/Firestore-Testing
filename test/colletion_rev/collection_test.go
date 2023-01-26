package collectionrev

import (
	"firebaseapi/helper"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_Collection(t *testing.T) {
	main_col := NewCollectionPayloads(WithRootCollection("hello"))
	col := main_col.Doc("hello").Col("test").Doc("this").Col("here")

	filters := []Filter{
		{
			By:  "cities",
			Op:  helper.EqualTo,
			Val: "USA",
		},
		{
			By:  "capital",
			Op:  helper.EqualTo,
			Val: false,
		},
	}

	sorts := []Sort{
		{
			By:  "cities",
			Dir: helper.ASC,
		},
	}

	for _, v := range filters {
		col = col.Where(v)
	}

	for _, v := range sorts {
		main_col = main_col.Sorts(v)
	}

	spew.Dump(col)
}
