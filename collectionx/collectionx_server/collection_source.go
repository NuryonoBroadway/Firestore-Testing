package collectionxserver

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/davecgh/go-spew/spew"
)

var (
	ErrInvalidPath = errors.New("invalid path")
)

// Collection Core Source Contract
type CollectionCore_SourceDocument interface {
	Retrive(ctx context.Context, p *Payload) (data *DataResponse, err error)
	Save(ctx context.Context, p *Payload) error
	Snapshots(ctx context.Context, p *Payload) (col *firestore.QuerySnapshotIterator, doc *firestore.DocumentSnapshotIterator, err error)
}

// Collection Core Source Implementation
type collectionCore_SourceDocumentImplementation struct {
	client *firestore.Client
	config *ServerConfig
}

func NewCollectionCore_SourceDocument(config *ServerConfig, client *firestore.Client) *collectionCore_SourceDocumentImplementation {
	return &collectionCore_SourceDocumentImplementation{client, config}
}

func (sd *collectionCore_SourceDocumentImplementation) Save(ctx context.Context, p *Payload) error {
	spew.Dump(p)

	if p == nil {
		return fmt.Errorf("payload is empty")
	}

	var (
		_, doc, _ = sd.pathBuilder(p)
	)

	return sd.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		data := make(map[string]interface{})
		for i := 0; i < len(p.Data.Row); i++ {
			data[p.Data.Row[i].Path] = p.Data.Row[i].Value
		}

		if p.Data.IsMergeAll {
			// universal set can create or update
			return tx.Set(doc, data, firestore.MergeAll)
		}

		return tx.Set(doc, data)

	})

}

func (sd *collectionCore_SourceDocumentImplementation) Retrive(ctx context.Context, p *Payload) (data *DataResponse, err error) {
	var (
		colRef, docRef, isLastDoc = sd.pathBuilder(p)
		isLastCol                 = !isLastDoc

		isFindAll = isLastCol
		isFindOne = isLastDoc
	)

	if isFindAll {
		docs, err := colRef.Data(p, ctx)
		if err != nil {
			return nil, err
		}

		size, err := colRef.Count(ctx)
		if err != nil {
			return nil, err
		}

		return &DataResponse{
			Size: size,
			Data: docs,
		}, nil
	} else if isFindOne {
		snap, err := docRef.Get(ctx)
		if err != nil {
			return nil, err
		}

		return &DataResponse{
			Size: 1,
			Data: snap.Data(),
		}, nil
	}

	return nil, ErrInvalidPath
}

func (sd *collectionCore_SourceDocumentImplementation) Snapshots(ctx context.Context, p *Payload) (col *firestore.QuerySnapshotIterator, doc *firestore.DocumentSnapshotIterator, err error) {
	var (
		colRef, docRef, isLastDoc = sd.pathBuilder(p)
		isLastCol                 = !isLastDoc

		isFindAll = isLastCol
		isFindOne = isLastDoc
	)

	if isFindAll {
		query := queryBuilder(colRef.colRef.Query, p)
		return query.Snapshots(ctx), nil, nil
	} else if isFindOne {
		return nil, docRef.Snapshots(ctx), nil
	}

	return nil, nil, errors.New("not implemented")
}
