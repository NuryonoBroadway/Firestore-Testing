package collectionxclient

import (
	"context"
	collectionxservice "firebaseapi/collectionx/collectionx_service"
	"firebaseapi/helper"
	"time"

	grpc "google.golang.org/grpc"
)

// Document cant search by filter
type Documenter interface {
	Col(id string) Collector
	Retrive() (*StandardAPI, error)
}

// Collector can search by filter
type Collector interface {
	Doc(id string) Documenter
	NewDoc() Documenter
	OrderBy(by string, dir OrderDir) *Payload
	Where(by string, op helper.Operator, val interface{}) *Payload
	Limit(limit int) *Payload
	DataRange(field string, start time.Time, end time.Time) *Payload
	Pagination(page int32, meta MetaData) *Payload
	Retrive() (*StandardAPI, error)
}

type Path struct {
	CollectionID string `json:"collection_id,omitempty"`
	DocumentID   string `json:"document_id,omitempty"`
	NewDocument  bool   `json:"new_document,omitempty"`
}

type Payload struct {
	conn *grpc.ClientConn
	ctx  context.Context

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

	// TODO: metadata
	query Filtering

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

func NewCollectionPayloads(opts ...func(p *Payload)) Documenter {
	p := Payload{}
	for _, v := range opts {
		v(&p)
	}

	return &p
}

func WithRootCollection(in string) func(p *Payload) {
	return func(p *Payload) {
		p.RootCollection = in
	}
}

func WithRootDocuments(in string) func(p *Payload) {
	if in == "" {
		in = "default"
	}

	return func(p *Payload) {
		p.RootDocument = in
	}
}

func WithGRPCCon(conn *grpc.ClientConn) func(p *Payload) {
	return func(p *Payload) {
		p.conn = conn
	}
}

func WithContext(ctx context.Context) func(p *Payload) {
	return func(p *Payload) {
		p.ctx = ctx
	}
}

type StandardAPI struct {
	Status  string `json:"status,omitempty"`
	Entity  string `json:"entity,omitempty"`
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
	Meta    Meta   `json:"meta,omitempty"`
	Data    Data   `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Meta struct {
	Page      int32 `json:"page"`
	PerPage   int32 `json:"per_page"`
	OrderBy   string
	OrderType string
}

type Data struct {
	Type string
	Data interface{}
}

type Error struct {
	General    string              `json:"general"`
	Validation []map[string]string `json:"validation"`
}

type ListValue struct {
	RefID  string                 `json:"ref_id"`
	Object map[string]interface{} `json:"object"`
}

func (p *Payload) Col(id string) Collector {
	p.Path = append(p.Path, Path{
		CollectionID: id,
	})
	return p
}

func (p *Payload) Doc(id string) Documenter {
	p.Path = append(p.Path, Path{
		DocumentID: id,
	})
	return p
}

func (p *Payload) NewDoc() Documenter {
	p.Path = append(p.Path, Path{
		NewDocument: true,
	})
	return p
}

func (p *Payload) Set(data map[string]interface{}) *Payload {
	p.Data = data
	return p
}

func (p *Payload) Limit(limit int) *Payload {
	p.limit = int32(limit)
	return p
}

func (p *Payload) OrderBy(by string, dir OrderDir) *Payload {
	p.query.Sort = Sort_Query{
		OrderBy:   by,
		OrderType: dir,
	}

	return p
}

func (p *Payload) DataRange(field string, start time.Time, end time.Time) *Payload {
	p.query.DateRange = DateRange_Query{
		Field: field,
		Start: start,
		End:   end,
	}

	return p
}

func (p *Payload) Where(by string, op helper.Operator, val interface{}) *Payload {
	filter := Filter_Query{
		By:  by,
		Op:  op.ToString(),
		Val: val,
	}

	p.query.Filter = append(p.query.Filter, filter)
	return p
}

func (p *Payload) Pagination(page int32, meta MetaData) *Payload {
	p.pagination = Pagination{
		Page: page,
	}

	if len(meta.Docs) == 0 {
		p.pagination.Meta = MetaData{
			Page: page,
			Docs: make([]map[string]interface{}, 0),
		}
	} else {
		p.pagination.Meta = meta
	}

	return p
}

func MetadataCreator(page int32, docs map[string]interface{}) MetaData {
	m := MetaData{
		Page: page,
	}

	for _, v := range docs {
		m.Docs = append(m.Docs, v.(map[string]interface{}))
	}

	return m
}

/*
Payload
*/

func (s *StandardAPI) MapValue() map[string]interface{} {
	switch s.Data.Type {
	case helper.Collection:
		if v, ok := s.Data.Data.([]ListValue); ok {
			mapped := map[string]interface{}{}
			for _, v := range v {
				mapped[v.RefID] = v.Object
			}

			return mapped
		}
	case helper.Document:
		if v, ok := s.Data.Data.(map[string]interface{}); ok {
			return v
		}
	}

	return nil
}
