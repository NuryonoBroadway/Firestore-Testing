package collectionxserver

import (
	"context"
	"errors"
	"firebaseapi/helper"

	"cloud.google.com/go/firestore"
)

var (
	ErrInvalidPath = errors.New("invalid path")
)

// Collection Core Source Contract
type CollectionCore_SourceDocument interface {
	Retrive(ctx context.Context, p *Payload) (data *DataResponse, err error)
	Snapshots(ctx context.Context, p *Payload) (col *firestore.QuerySnapshotIterator, doc *firestore.DocumentSnapshotIterator, err error)
}

// Collection Core Source Implementation
type collectionCore_SourceDocumentImplementation struct {
	client *firestore.Client
	config *ServerConfig
}

func NewCollectionCore_SourceDocument(config *ServerConfig) *collectionCore_SourceDocumentImplementation {
	client := RegistryFirestoreClient(config)
	return &collectionCore_SourceDocumentImplementation{client, config}
}

// Helper
func (sd *collectionCore_SourceDocumentImplementation) RootCollection(p *Payload) *firestore.CollectionRef {
	sd.config.RootCollection = p.RootCollection
	return sd.client.Collection(sd.config.RootCollection)
}

func (sd *collectionCore_SourceDocumentImplementation) RootDocument(p *Payload) *firestore.DocumentRef {
	sd.config.RootDocument = p.RootDocument
	return sd.RootCollection(p).Doc(sd.config.RootDocument)
}

func (sd *collectionCore_SourceDocumentImplementation) Builder(p *Payload) (*firestore.CollectionRef, *firestore.DocumentRef, bool) {
	var (
		isLastDoc bool
		docRef    *firestore.DocumentRef
		colRef    *firestore.CollectionRef
	)

	colRef = sd.RootCollection(p)
	docRef = sd.RootDocument(p)

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

func (sd *collectionCore_SourceDocumentImplementation) Retrive(ctx context.Context, p *Payload) (data *DataResponse, err error) {
	var (
		colRef, docRef, isLastDoc = sd.Builder(p)
		isLastCol                 = !isLastDoc

		isFindAll = isLastCol
		isFindOne = isLastDoc
	)

	// set default combination ordertype is asc and orderby is created_at, so if user want to use pagination
	// they dont need to specify ordertype and orderby
	// set the default getter to default orderby value

	if isFindAll {
		query := colRef.Query

		if p.query.Sort.OrderBy != "" {
			var dir firestore.Direction
			switch p.query.Sort.OrderType {
			case Asc:
				dir = firestore.Asc
			case Desc:
				dir = firestore.Desc
			}

			query = query.OrderBy(p.query.Sort.OrderBy, dir)
		}

		for i := 0; i < len(p.query.Filter); i++ {
			filter := p.query.Filter[i]
			query = query.Where(filter.By, filter.Op, filter.Val)
		}

		if p.query.DateRange.Field != "" {
			ranges := p.query.DateRange
			if ranges.Field == p.query.Sort.OrderBy {
				query = query.Where(
					p.query.DateRange.Field,
					helper.GreaterThanEqual.ToString(),
					p.query.DateRange.Start,
				).Where(
					p.query.DateRange.Field,
					helper.LessThanEqual.ToString(),
					p.query.DateRange.End,
				)
			} else {
				return nil, errors.New("use sort with the same field with daterange field")
			}
		}

		if p.limit > 0 {
			query = query.Limit(int(p.limit))
		}

		// limitation on firestore cursor is we cannot specify the page number we want
		if p.isPagination {
			page := p.pagination.Page
			offset := (page - 1) * p.limit
			query = query.Offset(int(offset))
		}

		snaps, err := query.Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}

		docs := make([]map[string]interface{}, len(snaps))
		for i := 0; i < len(snaps); i++ {
			docs[i] = map[string]interface{}{
				"ref_id": snaps[i].Ref.ID,
				"object": snaps[i].Data(),
			}
		}

		return &DataResponse{
			Total: len(snaps),
			Data:  docs,
		}, nil
	} else if isFindOne {
		snap, err := docRef.Get(ctx)
		if err != nil {
			return nil, err
		}

		return &DataResponse{
			Total: 1,
			Data:  snap.Data(),
		}, nil
	}

	return nil, ErrInvalidPath
}

func (sd *collectionCore_SourceDocumentImplementation) Snapshots(ctx context.Context, p *Payload) (col *firestore.QuerySnapshotIterator, doc *firestore.DocumentSnapshotIterator, err error) {
	var (
		colRef, docRef, isLastDoc = sd.Builder(p)
		isLastCol                 = !isLastDoc

		isFindAll = isLastCol
		isFindOne = isLastDoc
	)

	if isFindAll {
		query := colRef.Query
		if p.query.Sort.OrderBy != "" {
			var dir firestore.Direction
			switch p.query.Sort.OrderType {
			case Asc:
				dir = firestore.Asc
			case Desc:
				dir = firestore.Desc
			}

			query = query.OrderBy(p.query.Sort.OrderBy, dir)
		}

		for i := 0; i < len(p.query.Filter); i++ {
			filter := p.query.Filter[i]
			query = query.Where(filter.By, filter.Op, filter.Val)
		}

		if p.query.DateRange.Field != "" {
			ranges := p.query.DateRange
			if ranges.Field == p.query.Sort.OrderBy {
				query = query.Where(
					p.query.DateRange.Field,
					helper.GreaterThanEqual.ToString(),
					p.query.DateRange.Start,
				).Where(
					p.query.DateRange.Field,
					helper.LessThanEqual.ToString(),
					p.query.DateRange.End,
				)
			} else {
				return nil, nil, errors.New("use sort with the same field with daterange field")
			}
		}

		if p.limit > 0 {
			query = query.Limit(int(p.limit))
		}

		// limitation on firestore cursor is we cannot specify the page number we want
		if p.isPagination {
			page := p.pagination.Page
			offset := (page - 1) * p.limit
			query = query.Offset(int(offset))
		}

		return query.Snapshots(ctx), nil, nil
	} else if isFindOne {
		return nil, docRef.Snapshots(ctx), nil
	}

	return nil, nil, errors.New("not implemented")
}
