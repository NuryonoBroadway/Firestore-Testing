package collectionx

import (
	"context"
	"encoding/json"

	grpc "google.golang.org/grpc"
)

// Collection Core GRPC Server
func NewServer(source CollectionCore_SourceDocument, grpcOpt ...grpc.ServerOption) *grpc.Server {
	if source == nil {
		return nil
	}

	var (
		gsrv = grpc.NewServer(grpcOpt...)
		srv  = NewCollectionCoreServer(source)
	)

	RegisterServiceCollectionServer(gsrv, srv)

	return gsrv
}

// Collection Core Implementation
type server struct {
	source CollectionCore_SourceDocument
	UnimplementedServiceCollectionServer
}

func NewCollectionCoreServer(source CollectionCore_SourceDocument) *server {
	return &server{
		source: source,
	}
}

func (srv *server) Retrive(ctx context.Context, req *RetriveRequest) (*RetriveResponse, error) {
	if req.Payload == nil {
		return &RetriveResponse{
			Api: &StandardAPIProto{
				Status:  "ERROR",
				Entity:  "retriveFirestoreDocument",
				State:   "retriveFirestoreDocumentError",
				Message: "Retrive Firestore Document Failed Read Source Data",
			},
		}, nil
	}

	var (
		paths   = make([]Path, len(req.Payload.Path))
		sorts   = make([]Sort, len(req.Payload.Sort))
		filters = make([]Filter, len(req.Payload.Filter))
	)

	for i := 0; i < len(req.Payload.Path); i++ {
		paths[i].CollectionID = req.Payload.Path[i].CollectionId
		paths[i].DocumentID = req.Payload.Path[i].DocumentId
		paths[i].NewDocument = req.Payload.Path[i].NewDocument
	}

	for i := 0; i < len(req.Payload.Sort); i++ {
		sorts[i] = Sort{
			By:  req.Payload.Sort[i].By,
			Dir: req.Payload.Sort[i].Dir,
		}
	}

	for i := 0; i < len(req.Payload.Filter); i++ {
		filters[i] = Filter{
			By: req.Payload.Filter[i].By,
			Op: req.Payload.Filter[i].Op,
		}

		if req.Payload.Filter[i].GetValString() != "" {
			filters[i].Val = req.Payload.Filter[i].GetValString()
		} else if req.Payload.Filter[i].GetValInt() != 0 {
			filters[i].Val = req.Payload.Filter[i].GetValInt()
		} else {
			filters[i].Val = req.Payload.Filter[i].GetValBool()
		}
	}

	var (
		res = RetriveResponse{}
		p   = Payload{
			RootCollection: req.Payload.RootCollection,
			filter:         filters,
			limit:          req.Payload.Limit,
			sort:           sorts,
			IsDelete:       req.Payload.IsDelete,
			Data:           req.Payload.Data.AsMap(),
			Path:           paths,
		}
	)

	retriveRes, err := srv.source.Retrive(ctx, &p)
	if err != nil {
		res.Api = &StandardAPIProto{
			Status:  "ERROR",
			Entity:  "retriveFirestoreDocument",
			State:   "retriveFirestoreDocumentError",
			Message: "Retrive Firestore Document Failed Read Source Data",
			Error: &ErrorProto{
				General: err.Error(),
			},
		}
		return &res, nil
	}

	data, err := json.Marshal(retriveRes)
	if err != nil {
		res.Api = &StandardAPIProto{
			Status:  "ERROR",
			Entity:  "retriveFirestoreDocument",
			State:   "retriveFirestoreDocumentMarshalResponseError",
			Message: "Retrive Firestore Document Failed Build Result Data",
			Error: &ErrorProto{
				General: err.Error(),
			},
		}
		return &res, nil
	}

	res.Api = &StandardAPIProto{
		Status:  "SUCCESS",
		Entity:  "retriveFirestoreDocument",
		State:   "retriveFirestoreDocumentSuccess",
		Message: "Retrive Firestore Document Success",
	}
	res.Data = data

	return &res, nil
}
