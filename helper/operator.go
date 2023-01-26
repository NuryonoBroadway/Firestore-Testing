package helper

import (
	"reflect"
)

const (
	GreaterThan      = ">"
	LessThan         = "<"
	GreaterThanEqual = ">="
	LessThanEqual    = "<="
	EqualTo          = "=="
	NotEqualTo       = "!="
	NotIn            = "not-in"
	In               = "in"
	ArrayContains    = "array-contains"
	ArrayContainsAny = "array-contains-any"
)

const (
	ASC  = "asc"
	DESC = "desc"
)

const (
	Collection = "collection"
	Document   = "document"
)

func Compare(x interface{}, y interface{}) (int, int) {
	var xint int = 0
	var yint int = 0

	xtyp := reflect.TypeOf(x)
	switch xtyp.Kind() {
	case reflect.Int:
		xint = int(x.(int))
	case reflect.Int32:
		xint = int(x.(int32))
	case reflect.Int16:
		xint = int(x.(int16))
	case reflect.Int64:
		xint = int(x.(int64))
	}

	ytyp := reflect.TypeOf(y)
	switch ytyp.Kind() {
	case reflect.Int:
		yint = int(y.(int))
	case reflect.Int32:
		yint = int(y.(int32))
	case reflect.Int16:
		yint = int(y.(int16))
	case reflect.Int64:
		yint = int(y.(int64))
	}

	return xint, yint
}

func SliceCheckCondition(x interface{}, target interface{}) bool {
	m := make(map[interface{}]bool)
	if strings, ok := x.([]string); ok {
		for _, v := range strings {
			m[v] = true
		}

		return m[target]
	} else if ints, ok := x.([]int); ok {
		for _, v := range ints {
			m[v] = true
		}

		return m[target]
	}

	return false
}
