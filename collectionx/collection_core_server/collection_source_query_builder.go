package collection_core_server

import (
	"firebaseapi/helper"

	"cloud.google.com/go/firestore"
)

func (sd *sourceDocument) rootCollection(p *Payload) *firestore.CollectionRef {
	sd.config.RootCollection = p.RootCollection
	return sd.client.Collection(sd.config.RootCollection)
}

func (sd *sourceDocument) rootDocument(p *Payload) *firestore.DocumentRef {
	sd.config.RootDocument = p.RootDocument
	return sd.rootCollection(p).Doc(sd.config.RootDocument)
}

func (sd *sourceDocument) pathBuilder(p *Payload) (*firestore.CollectionRef, *firestore.DocumentRef, bool) {
	var (
		isLastDoc bool
		docRef    *firestore.DocumentRef
		colRef    *firestore.CollectionRef
	)

	colRef = sd.rootCollection(p)
	docRef = sd.rootDocument(p)

	// Path Builder
	for i := 0; i < len(p.Path); i++ {
		if p.Path[i].NewDocument {
			isLastDoc = true
			docRef = colRef.NewDoc()
		}

		if p.Path[i].DocumentID != "" {
			isLastDoc = true
			docRef = colRef.Doc(p.Path[i].DocumentID)
		}

		if p.Path[i].CollectionID != "" {
			isLastDoc = false
			colRef = docRef.Collection(p.Path[i].CollectionID)
		}
	}

	return colRef, docRef, isLastDoc
}

func (sd *sourceDocument) queryBuilder(query firestore.Query, p *Payload) firestore.Query {
	for i := 0; i < len(p.query.Sort); i++ {
		var dir firestore.Direction
		switch p.query.Sort[i].OrderType {
		case Asc:
			dir = firestore.Asc
		case Desc:
			dir = firestore.Desc
		}

		query = query.OrderBy(p.query.Sort[i].OrderBy, dir)
	}

	for i := 0; i < len(p.query.Filter); i++ {
		filter := p.query.Filter[i]
		query = query.Where(filter.By, filter.Op, filter.Val)
	}

	if len(p.query.DateRange.Field) > 0 {
		ranges := p.query.DateRange
		query = query.Where(
			ranges.Field,
			helper.GreaterThanEqual.ToString(),
			ranges.Start,
		).Where(
			ranges.Field,
			helper.LessThanEqual.ToString(),
			ranges.End,
		)
	}

	if p.limit > 0 {
		query = query.Limit(int(p.limit))
	}

	// limitation on firestore cursor is we cannot specify the page number we want
	if p.isPagination {
		page := p.page
		offset := (page - 1) * p.limit
		query = query.Offset(int(offset))
	}

	return query
}
