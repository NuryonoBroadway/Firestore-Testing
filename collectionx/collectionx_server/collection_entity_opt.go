package collectionxserver

import (
	collectionxservice "firebaseapi/collectionx/collectionx_service"
	"time"
)

type Path struct {
	CollectionID    string `json:"collection_id,omitempty"`
	DocumentID      string `json:"document_id,omitempty"`
	NewDocument     bool   `json:"new_document,omitempty"`
	CollectionGroup string `json:"collection_group,omitempty"`
}

type DataResponse struct {
	Size int
	Data interface{}
}

type Payload struct {
	RootCollection string `json:"root_collection"`
	RootDocument   string `json:"root_document"`

	// modify
	Data Datas `json:"data"`

	Path []Path `json:"path"`

	// limitation
	limit int32

	// pagination
	page int32

	// filtering
	query Filtering
	// condition
	isPagination bool
	isDelete     bool
}

type Datas struct {
	IsMergeAll bool  `json:"is_merge_all"`
	Row        []Row `json:"row"`
}

type Row struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type Pagination struct {
	Page int32
}

type Filtering struct {
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

type Aggregator struct {
	Name string `firestore:"name"`
	Size int    `firestore:"size"`
}

type ListValue struct {
	RefID  string                 `json:"ref_id"`
	Object map[string]interface{} `json:"object"`
}

type Message struct {
	ID        string
	Topic     string
	Data      []byte
	Attribute map[string]string
}

type Opts func(c *Options)

type Options struct {
	MaxConcurrent  int
	SubscribeAsync bool
	Topic          string
}

func defaults() *Options {
	return &Options{
		SubscribeAsync: false,
	}
}

func WithTopic(v string) Opts {
	return func(c *Options) {
		c.Topic = v
	}
}

func WithMaxConcurrent(v int) Opts {
	return func(c *Options) {
		c.MaxConcurrent = v
	}
}

func WithSubscribeAsync(v bool) Opts {
	return func(c *Options) {
		c.SubscribeAsync = v
	}
}
