package collectionrev

type Document struct {
	RootCollection string                 `json:"root_collection"`
	RootDocument   string                 `json:"root_document"`
	Data           map[string]interface{} `json:"data"`

	Path []Path `json:"path"`

	// condition
	IsDelete bool `json:"is_delete"`
	IsGet    bool `json:"is_get"`
}

func newDoc(parent *Collection, id string) *Document {
	parent.Path = append(parent.Path, Path{
		DocumentID: id,
	})
	return &Document{
		RootCollection: parent.RootCollection,
		Path:           parent.Path,
	}
}

func (d *Document) Col(id string) *Collection {
	return newCol(d, id)
}
