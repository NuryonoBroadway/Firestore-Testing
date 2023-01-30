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
	Snapshots() (*CollectionxSnapshots, error)
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
	Page(page int32) *Payload
	Snapshots() (*CollectionxSnapshots, error)
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

	// limitation
	limit int32

	// pagination
	pagination Pagination

	// filtering
	query Filtering

	// condition
	isPagination bool
	isDelete     bool
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
	p.pagination = Pagination{
		Page: page,
	}
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

/*
Document Change
*/

type Snapshots struct {
	StandardAPIDefault `json:"standard_api"`
	Kind               string
	Data               Data
	Timestamp          Timestamp
}

type CollectionxSnapshots struct {
	snapshots Snapshots

	isCol bool
	ws    collectionxservice.ServiceCollection_SnapshotsClient
	err   error
}

type DocumentKind int32

const (
	DOCUMENT_KIND_ADDED     DocumentKind = DocumentKind(collectionxservice.DocumentChangeKind_DOCUMENT_KIND_ADDED)
	DOCUMENT_KIND_REMOVED   DocumentKind = DocumentKind(collectionxservice.DocumentChangeKind_DOCUMENT_KIND_REMOVED)
	DOCUMENT_KIND_MODIFIED  DocumentKind = DocumentKind(collectionxservice.DocumentChangeKind_DOCUMENT_KIND_MODIFIED)
	DOCUMENT_KIND_SNAPSHOTS DocumentKind = DocumentKind(collectionxservice.DocumentChangeKind_DOCUMENT_KIND_SNAPSHOTS)
)

func (d DocumentKind) ToString() string {
	switch d {
	case DOCUMENT_KIND_ADDED:
		return "Document Added"
	case DOCUMENT_KIND_REMOVED:
		return "Document Removed"
	case DOCUMENT_KIND_MODIFIED:
		return "Document Modified"
	case DOCUMENT_KIND_SNAPSHOTS:
		return "Document Snapshots"
	default:
		return "not-implemented"
	}
}

type Timestamp struct {
	CreatedTime time.Time `json:"created_time"`
	ReadTime    time.Time `json:"read_time"`
	UpdateTime  time.Time `json:"modified_time"`
}
