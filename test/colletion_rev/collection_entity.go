package collectionrev

import "firebaseapi/helper"

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
	By  string          `json:"by"`
	Op  helper.Operator `json:"op"`
	Val interface{}     `json:"val"`
}
