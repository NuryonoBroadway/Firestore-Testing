package collectionxclient

import (
	collectionxservice "firebaseapi/collectionx/collectionx_service"
	"firebaseapi/helper"
	"time"
)

/*
Standar API
*/

type StandardAPIDefault struct {
	Status         string          `json:"status,omitempty"`
	Entity         string          `json:"entity,omitempty"`
	State          string          `json:"state,omitempty"`
	Representative *Representative `json:"representative,omitempty"`
}

type StandardAPI struct {
	*StandardAPIDefault `json:"standard_api"`
	Meta                Meta `json:"meta,omitempty"`
	Data                Data `json:"data,omitempty"`
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

type Representative struct {
	Code    helper.Code `json:"code"`
	Message string      `json:"message"`
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
	*StandardAPIDefault `json:"standard_api"`
	Kind                string
	Data                Data
	Timestamp           Timestamp
}

type CollectionxSnapshots struct {
	res *collectionxservice.SnapshotsResponse

	client *client
	isCol  bool
	ws     collectionxservice.ServiceCollection_SnapshotsClient
	err    error
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

func NewSuccessResponse() *StandardAPIDefault {
	return &StandardAPIDefault{
		Status: "SUCCESS",
		Entity: "retrieveFirestore",
		State:  "retrieveFirestoreSuccess",
	}
}

func NewErrorResponse() *StandardAPIDefault {
	return &StandardAPIDefault{
		Status: "ERROR",
		Entity: "retriveFirestore",
		State:  "retrieveFirestoreError",
	}
}

func (sa *StandardAPIDefault) WithRepresentative(code helper.Code, message string) *StandardAPIDefault {
	sa.Representative = &Representative{
		Code:    code,
		Message: message,
	}
	return sa
}

