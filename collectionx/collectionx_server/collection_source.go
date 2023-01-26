package collectionxserver

import (
	"context"
	"errors"
	"firebaseapi/helper"
	"strings"

	"cloud.google.com/go/firestore"
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
	return sd.client.Collection(p.RootCollection)
}

func (sd *collectionCore_SourceDocumentImplementation) RootDocument(p *Payload) *firestore.DocumentRef {
	return sd.RootCollection(p).Doc(p.Path[0].DocumentID)
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

	if isFindAll {
		q := colRef.Query

		if p.sort.Dir != "" {
			var dir firestore.Direction
			switch strings.ToLower(p.sort.Dir) {
			case helper.ASC:
				dir = firestore.Asc
			case helper.DESC:
				dir = firestore.Desc
			}

			q = q.OrderBy(p.sort.By, dir)
		}

		if len(p.filter) != 0 {
			for i := 0; i < len(p.filter); i++ {
				filter := p.filter[i]
				q = q.Where(filter.By, filter.Op, filter.Val)
			}
		}

		if p.limit > 0 {
			q = q.Limit(int(p.limit))
		}

		snaps, err := q.Documents(ctx).GetAll()
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
