package collectionxserver

import (
	"context"
	"encoding/json"
	collectionxservice "firebaseapi/collectionx/collectionx_service"

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

	collectionxservice.RegisterServiceCollectionServer(gsrv, srv)

	return gsrv
}

// Collection Core Implementation
type server struct {
	source CollectionCore_SourceDocument
	collectionxservice.UnimplementedServiceCollectionServer
}

func NewCollectionCoreServer(source CollectionCore_SourceDocument) *server {
	return &server{
		source: source,
	}
}

func (srv *server) Retrive(ctx context.Context, req *collectionxservice.RetriveRequest) (*collectionxservice.RetriveResponse, error) {
	if req.Payload == nil {
		return &collectionxservice.RetriveResponse{
			Api: &collectionxservice.StandardAPIProto{
				Status:  "ERROR",
				Entity:  "retriveFirestoreDocument",
				State:   "retriveFirestoreDocumentError",
				Message: "Retrive Firestore Document Failed Read Source Data",
			},
		}, nil
	}

	var (
		paths      = make([]Path, len(req.Payload.Path))
		query      = Filtering{}
		pagination = Pagination{}
	)

	for i := 0; i < len(req.Payload.Path); i++ {
		paths[i].CollectionID = req.Payload.Path[i].CollectionId
		paths[i].DocumentID = req.Payload.Path[i].DocumentId
		paths[i].NewDocument = req.Payload.Path[i].NewDocument
	}

	query.Sort = Sort_Query{
		OrderBy:   req.Payload.Query.Sort.OrderBy,
		OrderType: OrderDir(req.Payload.Query.Sort.OrderType.Number()),
	}

	query.DateRange = DateRange_Query{
		Field: req.Payload.Query.DateRange.Field,
		Start: req.GetPayload().Query.DateRange.Start.AsTime(),
		End:   req.GetPayload().Query.DateRange.End.AsTime(),
	}

	filters := make([]Filter_Query, len(req.Payload.Query.Filter))
	for i := 0; i < len(req.Payload.Query.Filter); i++ {
		filters[i] = Filter_Query{
			By: req.Payload.Query.Filter[i].By,
			Op: req.Payload.Query.Filter[i].Op,
		}
		if req.Payload.Query.Filter[i].GetValString() != "" {
			filters[i].Val = req.Payload.Query.Filter[i].GetValString()
		} else if req.Payload.Query.Filter[i].GetValInt() < -1 {
			filters[i].Val = req.Payload.Query.Filter[i].GetValInt()
		} else {
			filters[i].Val = req.Payload.Query.Filter[i].GetValBool()
		}
	}
	query.Filter = filters

	if req.Payload.Pagination.Meta != nil {
		pagination = Pagination{
			Page: req.Payload.Pagination.Page,
			Meta: MetaData{
				Page: req.Payload.Pagination.Meta.Page,
			},
		}

		var page []map[string]interface{}
		if err := json.Unmarshal(req.Payload.Pagination.Meta.Docs, &page); err != nil {
			return nil, err
		}
		pagination.Meta.Docs = page
	}

	var (
		res = collectionxservice.RetriveResponse{}
		p   = Payload{
			RootCollection: req.Payload.RootCollection,
			RootDocument:   req.Payload.RootDocument,
			limit:          req.Payload.Limit,
			IsDelete:       req.Payload.IsDelete,
			Data:           req.Payload.Data.AsMap(),
			Path:           paths,
			pagination:     pagination,
			query:          query,
		}
	)

	retriveRes, err := srv.source.Retrive(ctx, &p)
	if err != nil {
		res.Api = &collectionxservice.StandardAPIProto{
			Status:  "ERROR",
			Entity:  "retriveFirestoreDocument",
			State:   "retriveFirestoreDocumentError",
			Message: "Retrive Firestore Document Failed Read Source Data",
			Error: &collectionxservice.ErrorProto{
				General: err.Error(),
			},
		}
		return &res, nil
	}

	data, err := json.Marshal(retriveRes)
	if err != nil {
		res.Api = &collectionxservice.StandardAPIProto{
			Status:  "ERROR",
			Entity:  "retriveFirestoreDocument",
			State:   "retriveFirestoreDocumentMarshalResponseError",
			Message: "Retrive Firestore Document Failed Build Result Data",
			Error: &collectionxservice.ErrorProto{
				General: err.Error(),
			},
		}
		return &res, nil
	}

	res.Api = &collectionxservice.StandardAPIProto{
		Status:  "SUCCESS",
		Entity:  "retriveFirestoreDocument",
		State:   "retriveFirestoreDocumentSuccess",
		Message: "Retrive Firestore Document Success",
		Meta: &collectionxservice.MetaProto{
			Page:      p.pagination.Page,
			PerPage:   p.limit,
			OrderBy:   p.query.Sort.OrderBy,
			OrderType: p.query.Sort.OrderType.ToString(),
		},
	}
	res.Data = data

	return &res, nil
}
