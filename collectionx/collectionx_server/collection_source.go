package collectionxserver

import (
	"context"
	"errors"

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
func (sd *collectionCore_SourceDocumentImplementation) rootCollection(p *Payload) *firestore.CollectionRef {
	sd.config.RootCollection = p.RootCollection
	return sd.client.Collection(sd.config.RootCollection)
}

func (sd *collectionCore_SourceDocumentImplementation) rootDocument(p *Payload) *firestore.DocumentRef {
	sd.config.RootDocument = p.RootDocument
	return sd.rootCollection(p).Doc(sd.config.RootDocument)
}

func (sd *collectionCore_SourceDocumentImplementation) Retrive(ctx context.Context, p *Payload) (data *DataResponse, err error) {
	var (
		colRef, docRef, isLastDoc = pathBuilder(sd, p)
		isLastCol                 = !isLastDoc

		isFindAll = isLastCol
		isFindOne = isLastDoc
	)

	if isFindAll {
		query := queryBuilder(colRef.Query, p)
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
		colRef, docRef, isLastDoc = pathBuilder(sd, p)
		isLastCol                 = !isLastDoc

		isFindAll = isLastCol
		isFindOne = isLastDoc
	)

	if isFindAll {
		query := queryBuilder(colRef.Query, p)
		return query.Snapshots(ctx), nil, nil
	} else if isFindOne {
		return nil, docRef.Snapshots(ctx), nil
	}

	return nil, nil, errors.New("not implemented")
}
