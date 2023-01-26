package collectionxserver

type Sort struct {
	By  string `json:"by"`
	Dir string `json:"dir"`
}

type Filter struct {
	By  string      `json:"by"`
	Op  string      `json:"op"`
	Val interface{} `json:"val"`
}

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
