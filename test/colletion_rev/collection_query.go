package collectionrev

type Query struct {
	// pagination
	Limit int32 `json:"limit"`

	// TODO: filtering
	Filter []Filter `json:"filter"`

	// TODO: sorting
	Sort []Sort `json:"sort"`
}

func (c *Collection) Where(filter Filter) *Collection {
	c.Query.Filter = append(c.Query.Filter, filter)
	return c
}

func (c *Collection) Limits(how int) *Collection {
	c.Query.Limit = int32(how)
	return c
}

func (c *Collection) Sorts(sort Sort) *Collection {
	c.Query.Sort = append(c.Query.Sort, sort)
	return c
}
