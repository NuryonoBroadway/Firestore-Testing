package collectionx

import (
	"context"
	"errors"
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
	config *Config
}

func NewCollectionCore_SourceDocument(config *Config) *collectionCore_SourceDocumentImplementation {
	client := registryFirestoreClient(*config)
	return &collectionCore_SourceDocumentImplementation{client, config}
}

// Helper
func (sd *collectionCore_SourceDocumentImplementation) RootCollection() *firestore.CollectionRef {
	return sd.client.Collection(sd.config.ExternalCollection)
}

func (sd *collectionCore_SourceDocumentImplementation) RootDocument(root *firestore.CollectionRef) *firestore.DocumentRef {
	return root.Doc(sd.config.ExternalDocument)
}

func (sd *collectionCore_SourceDocumentImplementation) Builder(p *Payload) (*firestore.CollectionRef, *firestore.DocumentRef, bool) {
	var (
		isLastDoc bool
		docRef    *firestore.DocumentRef
		colRef    *firestore.CollectionRef

		// rootCollection = ""
	)

	// if p.RootCollection != "" {
	// 	rootCollection = p.RootCollection
	// }

	colRef = sd.RootCollection()
	docRef = sd.RootDocument(colRef)

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
		for i := 0; i < len(p.sort); i++ {
			var dir firestore.Direction
			switch strings.ToLower(p.sort[i].Dir) {
			case "asc":
				dir = firestore.Asc
			case "desc":
				dir = firestore.Desc
			}

			q = q.OrderBy(p.sort[i].By, dir)
		}

		for i := 0; i < len(p.filter); i++ {
			filter := p.filter[i]
			q = q.Where(filter.By, filter.Op, filter.Val)
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
