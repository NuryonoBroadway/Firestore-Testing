package collection_core_server

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
)

var (
	ErrInvalidPath = errors.New("invalid query path")
)

type SourceDocument interface {
	Retrive(ctx context.Context, p *Payload) (data *SourceData, err error)
}

type sourceDocument struct {
	client *firestore.Client
	config *ServerConfig
}

func NewSourceDocument(config *ServerConfig) *sourceDocument {
	client := RegistryFirestoreClient(config)
	return &sourceDocument{client, config}
}

func (sd *sourceDocument) Retrive(ctx context.Context, p *Payload) (data *SourceData, err error) {
	var (
		colRef, docRef, isLastDoc = sd.pathBuilder(p)
		isLastCol                 = !isLastDoc
		isFindAll                 = isLastCol
		isFindOne                 = isLastDoc
	)

	if isFindAll {
		var (
			query      = sd.queryBuilder(colRef.Query, p)
			snaps, err = query.Documents(ctx).GetAll()
		)
		if err != nil {
			return nil, err
		}

		docs := make([]ObjectData, len(snaps))
		for i := 0; i < len(snaps); i++ {
			docs[i] = ObjectData{
				RefID:  snaps[i].Ref.ID,
				Object: snaps[i].Data(),
			}
		}

		return &SourceData{
			Size: len(snaps),
			Data: docs,
		}, nil
	} else if isFindOne {
		snap, err := docRef.Get(ctx)
		if err != nil {
			return nil, err
		}

		return &SourceData{
			Size: 1,
			Data: snap.Data(),
		}, nil
	}

	return nil, ErrInvalidPath
}
