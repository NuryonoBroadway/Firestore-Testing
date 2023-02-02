package collectionxserver

import (
	"context"
	"errors"
	"firebaseapi/helper"
	"fmt"

	"cloud.google.com/go/firestore"
)

func (sd *collectionCore_SourceDocumentImplementation) pathBuilder(p *Payload) (*collectionAggregator, *firestore.DocumentRef, bool) {
	var (
		isLastDoc   bool
		docRef      *firestore.DocumentRef
		colRef      *firestore.CollectionRef
		colGroupRef *firestore.CollectionGroupRef
	)

	// Path Builder
	for i := 0; i < len(p.Path); i++ {
		if p.Path[i].NewDocument {
			isLastDoc = true
			docRef = colRef.NewDoc()
		}

		if p.Path[i].CollectionGroup != "" {
			isLastDoc = false
			colGroupRef = sd.client.CollectionGroup(p.Path[i].CollectionGroup)
		}

		if p.Path[i].DocumentID != "" {
			isLastDoc = true

			if i == 0 {
				docRef = sd.client.Doc(p.Path[i].DocumentID)
			} else {
				docRef = colRef.Doc(p.Path[i].DocumentID)
			}
		}

		if p.Path[i].CollectionID != "" {
			isLastDoc = false

			if i == 0 {
				colRef = sd.client.Collection(p.Path[i].CollectionID)
			} else {
				colRef = docRef.Collection(p.Path[i].CollectionID)
			}
		}
	}

	return &collectionAggregator{
		colRef:      colRef,
		colGroupRef: colGroupRef,
	}, docRef, isLastDoc
}

func queryBuilder(query firestore.Query, p *Payload) firestore.Query {
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

	if p.query.DateRange.Field != "" {
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

type collectionAggregator struct {
	colRef      *firestore.CollectionRef
	colGroupRef *firestore.CollectionGroupRef
}

func (ca *collectionAggregator) Data(p *Payload, ctx context.Context) ([]ListValue, error) {
	var query firestore.Query
	if ca.colGroupRef != nil {
		query = queryBuilder(ca.colGroupRef.Query, p)
	} else if ca.colRef != nil {
		query = queryBuilder(ca.colRef.Query, p)
	}

	snaps, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	docs := make([]ListValue, len(snaps))
	for i := 0; i < len(snaps); i++ {
		docs[i] = ListValue{
			RefID:  snaps[i].Ref.ID,
			Object: snaps[i].Data(),
		}
	}

	return docs, nil
}

func (ca *collectionAggregator) Count(ctx context.Context) (int, error) {
	if ca.colRef != nil {
		builder := fmt.Sprintf("aggregator_%v", ca.colRef.ID)
		total, err := ca.colRef.Doc(builder).Get(ctx)
		if err != nil {
			return 0, err
		}

		var a Aggregator
		if err := total.DataTo(&a); err != nil {
			return 0, err
		}

		return a.Size, nil
	} else if ca.colGroupRef != nil {
		return 0, nil
	}

	return 0, errors.New("not-implemented")
}
