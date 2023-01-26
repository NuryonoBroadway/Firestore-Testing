package collectionrev

type Collection struct {
	RootCollection string                 `json:"root_collection"`
	RootDocument   string                 `json:"root_document"`
	Data           map[string]interface{} `json:"data"`

	Path []Path `json:"path"`

	Query

	// condition
	IsDelete bool `json:"is_delete"`
	IsGet    bool `json:"is_get"`
}

func newCol(parent *Document, id string) *Collection {
	parent.Path = append(parent.Path, Path{
		CollectionID: id,
	})
	return &Collection{
		RootCollection: parent.RootCollection,
		Path:           parent.Path,
	}
}

func (c *Collection) Doc(id string) *Document {
	return newDoc(c, id)
}

type CollectionPayloadOption func(p *Collection)

func NewCollectionPayloads(opts ...CollectionPayloadOption) *Collection {
	p := Collection{}
	for _, v := range opts {
		v(&p)
	}

	return &p
}

// This option, will be replace other root collection name id
// if this option value not empty
func WithRootCollection(in string) CollectionPayloadOption {
	return func(p *Collection) {
		p.RootCollection = in
	}
}
