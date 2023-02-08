package collectionxclient

import (
	collectionxservice "firebaseapi/collectionx/collectionx_service"
	"firebaseapi/helper"
	"time"
)

// Document cant search by filter
type Documenter interface {
	Col(id string) Collector
	Set(data []Row, merge bool) Modify
	Retrive() (*StandardAPI, error)
	Snapshots() (*CollectionxSnapshots, error)
}

type Modify interface {
	Save() (*StandardAPIDefault, error)
}

// Collector can search by filter
type Collector interface {
	Doc(id string) Documenter
	NewDoc() Documenter
	OrderBy(by string, dir OrderDir) *Payload
	Where(by string, op helper.Operator, val interface{}) *Payload
	Limit(limit int) *Payload
	DateRange(field string, start time.Time, end time.Time) *Payload
	Page(page int32) *Payload
	Snapshots() (*CollectionxSnapshots, error)
	Retrive() (*StandardAPI, error)
}

type CollectorGroup interface {
	OrderBy(by string, dir OrderDir) *Payload
	Where(by string, op helper.Operator, val interface{}) *Payload
	Limit(limit int) *Payload
	DateRange(field string, start time.Time, end time.Time) *Payload
	Page(page int32) *Payload
	Snapshots() (*CollectionxSnapshots, error)
	Retrive() (*StandardAPI, error)
}

type Path struct {
	CollectionID    string `json:"collection_id,omitempty"`
	DocumentID      string `json:"document_id,omitempty"`
	NewDocument     bool   `json:"new_document,omitempty"`
	CollectionGroup string `json:"collection_group,omitempty"`
}

type Payload struct {
	client *client

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

func (p *Payload) Limit(limit int) *Payload {
	p.limit = int32(limit)
	return p
}

func (p *Payload) OrderBy(by string, dir OrderDir) *Payload {
	for i := 0; i < len(p.query.Sort); i++ {
		if p.query.Sort[i].OrderBy == by {
			return p
		}
	}

	sort := Sort_Query{
		OrderBy:   by,
		OrderType: dir,
	}

	p.query.Sort = append(p.query.Sort, sort)
	return p
}

func (p *Payload) DateRange(field string, start time.Time, end time.Time) *Payload {
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

func (p *Payload) Set(data []Row, merge bool) Modify {
	p.Data = Datas{
		IsMergeAll: merge,
		Row:        data,
	}

	return p
}
