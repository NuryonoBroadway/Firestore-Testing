package collectionxclient

import (
	"context"
	"firebaseapi/collectionx/collection_core_service"
	"firebaseapi/helper"
	"time"

	grpc "google.golang.org/grpc"
)

type Documenter interface {
	Col(id string) Collector
	Retrive() (*StandardAPI, error)
}

type Collector interface {
	Doc(id string) Documenter
	NewDoc() Documenter
	OrderBy(by string, dir OrderDir) *Payload
	Where(by string, op helper.Operator, val interface{}) *Payload
	Limit(limit int) *Payload
	DataRange(field string, start time.Time, end time.Time) *Payload
	Page(page int32) *Payload
	Retrive() (*StandardAPI, error)
}

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
	conn *grpc.ClientConn
	ctx  context.Context

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

type OptionPayload func(p *Payload)

func NewCollectionPayloads(opts ...OptionPayload) Documenter {
	p := Payload{}
	for i := 0; i < len(opts); i++ {
		opts[i](&p)
	}
	return &p
}

func WithRootCollection(in string) OptionPayload {
	return func(p *Payload) {
		p.RootCollection = in
	}
}

func WithRootDocuments(in string) OptionPayload {
	if in == "" {
		in = "default"
	}

	return func(p *Payload) {
		p.RootDocument = in
	}
}

func WithGRPCCon(conn *grpc.ClientConn) OptionPayload {
	return func(p *Payload) {
		p.conn = conn
	}
}

func WithContext(ctx context.Context) OptionPayload {
	return func(p *Payload) {
		p.ctx = ctx
	}
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
	sort := Sort_Query{
		OrderBy:   by,
		OrderType: dir,
	}

	p.query.Sort = append(p.query.Sort, sort)
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

func (p *Payload) Page(page int32) *Payload {
	p.isPagination = true
	p.page = page
	return p
}

/*
	Standar API
*/

type StandardAPIDefault struct {
	Status  string `json:"status,omitempty"`
	Entity  string `json:"entity,omitempty"`
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type StandardAPI struct {
	StandardAPIDefault `json:"standard_api"`
	Meta               Meta `json:"meta,omitempty"`
	Data               Data `json:"data,omitempty"`
}

type Meta struct {
	Page    int32 `json:"page"`
	PerPage int32 `json:"per_page"`
	Size    int32 `json:"size"`
	Total   int32 `json:"total"`
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
