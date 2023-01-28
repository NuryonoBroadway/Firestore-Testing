package collectionxserver

import (
	"context"
	"errors"
	"firebaseapi/helper"

	"cloud.google.com/go/firestore"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidPath = errors.New("invalid path")
)

// Collection Core Source Contract
type CollectionCore_SourceDocument interface {
	Retrive(ctx context.Context, p *Payload) (data interface{}, err error)
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

func (sd *collectionCore_SourceDocumentImplementation) Retrive(ctx context.Context, p *Payload) (data interface{}, err error) {
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

		if len(p.pagination.Meta.Docs) != 0 {
			// find last index
			// if let say we gonna get the meta for the current page use have
			current_page := p.pagination.Meta.Page
			page := p.pagination.Page
			if current_page == page {
				logrus.Info("in same page")
				query = query.StartAt(p.pagination.Meta.Docs[0][p.query.Sort.OrderBy]).Limit(int(p.limit))
			} else if current_page > page {
				logrus.Info("previous page")
				query = query.EndAt(p.pagination.Meta.Docs[0][p.query.Sort.OrderBy]).LimitToLast(int(p.limit))
			} else if current_page < page {
				logrus.Info("next page")
				data := p.pagination.Meta.Docs
				query = query.StartAfter(data[len(data)-1][p.query.Sort.OrderBy]).Limit(int(p.limit))
			}

		} else {
			if p.limit > 0 {
				query = query.Limit(int(p.limit))
			}
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

		return docs, nil
	} else if isFindOne {
		snap, err := docRef.Get(ctx)
		if err != nil {
			return nil, err
		}

		return snap.Data(), nil
	}

	return nil, ErrInvalidPath
}
