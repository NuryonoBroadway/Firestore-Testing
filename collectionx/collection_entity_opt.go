package collectionx

import (
	"context"

	grpc "google.golang.org/grpc"
)

type Setter interface {
	Set(data map[string]interface{}) *Payload
}

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
	Order(sort ...Sort) *Payload
	Max(limit int) *Payload
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

type ClientFir struct {
	Environment    string `json:"environment"`
	ServiceName    string `json:"service_name"`
	ProjectName    string `json:"project_name"`
	RootCollection string `json:"root_collection"`
}

type Payload struct {
	conn *grpc.ClientConn
	ctx  context.Context

	RootCollection string                 `json:"root_collection"`
	RootDocument   string                 `json:"root_document"`
	Data           map[string]interface{} `json:"data"`
	Path           []Path                 `json:"path"`

	// pagination
	limit int32

	// TODO: filtering
	filter []Filter

	// TODO: sorting
	sort []Sort

	// condition
	IsDelete bool `json:"is_delete"`
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

func (p *Payload) Max(limit int) *Payload {
	p.limit = int32(limit)
	return p
}

func (p *Payload) Order(sort ...Sort) *Payload {
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

type CollectionPayloadOption func(p *Payload)

func NewCollectionPayloads(opts ...CollectionPayloadOption) Collector {
	p := Payload{}
	for _, v := range opts {
		v(&p)
	}

	return &p
}

// This option, will be replace other root collection name id
// if this option value not empty
func WithRootCollection(in string) CollectionPayloadOption {
	return func(p *Payload) {
		p.RootCollection = in
	}
}

func WithGRPCCon(conn *grpc.ClientConn) CollectionPayloadOption {
	return func(p *Payload) {
		p.conn = conn
	}
}

func WithContext(ctx context.Context) CollectionPayloadOption {
	return func(p *Payload) {
		p.ctx = ctx
	}
}
