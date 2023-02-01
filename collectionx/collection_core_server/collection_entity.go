package collection_core_server

import (
	"firebaseapi/collectionx/collection_core_service"
	"time"
)

type Path struct {
	CollectionID string `json:"collection_id,omitempty"`
	DocumentID   string `json:"document_id,omitempty"`
	NewDocument  bool   `json:"new_document,omitempty"`
}

type SourceData struct {
	Size int
	Data interface{}
}

type ObjectData struct {
	RefID  string                 `json:"ref_id"`
	Object map[string]interface{} `json:"object"`
}

type Payload struct {
	RootCollection string
	RootDocument   string

	Data map[string]interface{}
	Path []Path

	// pagination
	limit int32
	page  int32
	query Queries

	// condition
	isPagination bool
	isDelete     bool
}

// Queries

type Queries struct {
	Sort      []Sort_Query
	Filter    []Filter_Query
	DateRange DateRange_Query
}

type Sort_Query struct {
	OrderBy   string
	OrderType OrderDir
}

type Filter_Query struct {
	By  string
	Op  string
	Val interface{}
}

type DateRange_Query struct {
	Field string
	Start time.Time
	End   time.Time
}

// Sorting

type OrderDir int32

const (
	Asc  OrderDir = OrderDir(collection_core_service.OrderTypeProto_ORDER_TYPE_ASC)
	Desc OrderDir = OrderDir(collection_core_service.OrderTypeProto_ORDER_TYPE_DESC)
)

func (o OrderDir) ToString() string {
	switch o {
	case Asc:
		return "ASC"
	case Desc:
		return "DESC"
	default:
		return ""
	}
}
