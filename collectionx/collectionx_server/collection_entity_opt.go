package collectionxserver

import (
	collectionxservice "firebaseapi/collectionx/collectionx_service"
	"time"
)

type Path struct {
	CollectionID string `json:"collection_id,omitempty"`
	DocumentID   string `json:"document_id,omitempty"`
	NewDocument  bool   `json:"new_document,omitempty"`
}

type Payload struct {
	Environment    string
	ServiceName    string
	ProjectName    string
	RootCollection string
	RootDocument   string

	Data map[string]interface{}
	Path []Path

	// pagination
	limit int32

	pagination Pagination
	query      Filtering

	// condition
	IsDelete bool
}

type Pagination struct {
	Page int32
	Meta MetaData
}

type MetaData struct {
	Page int32
	Docs []map[string]interface{}
}

type Filtering struct {
	Sort      Sort_Query
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

type OrderDir int32

const (
	Asc  OrderDir = OrderDir(collectionxservice.OrderTypeProto_ORDER_TYPE_ASC)
	Desc OrderDir = OrderDir(collectionxservice.OrderTypeProto_ORDER_TYPE_DESC)
)

func (o OrderDir) ToString() string {
	switch o {
	case Asc:
		return "ASC"
	case Desc:
		return "DESC"
	default:
		return "not-implemented"
	}
}
