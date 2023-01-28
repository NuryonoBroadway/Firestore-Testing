package helper_test

import (
	"firebaseapi/helper"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestSliceCheckCondition(t *testing.T) {
	type filter struct {
		Operator helper.Operator
		Value    interface{}
	}

	type data struct {
		name  string
		value map[string]interface{}
	}

	value := []data{
		{
			name: "WD",
			value: map[string]interface{}{
				"name":    "Washington D.C.",
				"country": "USA",
			},
		},
		{
			name: "SF",
			value: map[string]interface{}{
				"name":    "San Francisco",
				"country": "USA",
			},
		},
		{
			name: "LA",
			value: map[string]interface{}{
				"name":    "Los Angeles",
				"country": "USA",
			},
		},
		{
			name: "TK",
			value: map[string]interface{}{
				"name":    "Tokyo",
				"country": "Japan",
			},
		},
		{
			name: "BJ",
			value: map[string]interface{}{
				"name":    "Beijing",
				"country": "China",
			},
		},
	}

	slice := filter{
		Operator: helper.NotIn,
		Value:    []string{"China", "Japan"},
	}

	var newData []data
	for _, v := range value {
		switch slice.Operator {
		case helper.In:
			if helper.SliceCheckCondition(slice.Value, v.value["country"]) {
				data := data{
					name:  v.name,
					value: v.value,
				}

				newData = append(newData, data)
			}

		case helper.NotIn:
			if !helper.SliceCheckCondition(slice.Value, v.value["country"]) {
				data := data{
					name:  v.name,
					value: v.value,
				}

				newData = append(newData, data)
			}
		}
	}

	spew.Dump(newData)
}
