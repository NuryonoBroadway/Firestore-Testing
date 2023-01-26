package collectionxclient

import (
	"context"
	"firebaseapi/helper"

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
	Where(filter ...Filter) *Payload
	Order(sort Sort) *Payload
	Limit(limit int) *Payload
	Retrive() (*StandardAPI, error)
}

type Path struct {
	CollectionID string `json:"collection_id,omitempty"`
	DocumentID   string `json:"document_id,omitempty"`
	NewDocument  bool   `json:"new_document,omitempty"`
}

type Sort struct {
	By  string `json:"by"`
	Dir string `json:"dir"`
}

type Filter struct {
	By  string      `json:"by"`
	Op  string      `json:"op"`
	Val interface{} `json:"val"`
}

type Payload struct {
	conn *grpc.ClientConn
	ctx  context.Context

	Environment    string
	ServiceName    string
	ProjectName    string
	RootCollection string
	Data           map[string]interface{}
	Path           []Path

	// pagination
	limit int32

	// TODO: filtering
	filter []Filter

	// TODO: sorting
	sort Sort

	// condition
	IsDelete bool
}

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
		p.Path = append(p.Path, Path{
			DocumentID: in,
		})
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
	Data    Data   `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
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

func (p *Payload) Order(sort Sort) *Payload {
	p.sort = sort
	return p
}

func (p *Payload) Where(filter ...Filter) *Payload {
	p.filter = filter
	return p
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
