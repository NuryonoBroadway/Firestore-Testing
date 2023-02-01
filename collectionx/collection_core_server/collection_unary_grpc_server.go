package collection_core_server

import (
	"context"
	"encoding/json"
	"firebaseapi/collectionx/collection_core_service"

	grpc "google.golang.org/grpc"
)

// Collection Core GRPC Server
func NewServer(source SourceDocument, grpcOpt ...grpc.ServerOption) *grpc.Server {
	if source == nil {
		return nil
	}

	var (
		gsrv = grpc.NewServer(grpcOpt...)
		srv  = NewCollectionCoreServer(source)
	)

	collection_core_service.RegisterServiceCollectionServer(gsrv, srv)

	return gsrv
}

// Collection Core Implementation
type server struct {
	source SourceDocument
	collection_core_service.UnimplementedServiceCollectionServer
}

func NewCollectionCoreServer(source SourceDocument) *server {
	return &server{
		source: source,
	}
}

func (srv *server) payloadBuilder(req *collection_core_service.PayloadProto) (path_query []Path, filter_query Queries, page int32, err error) {
	var (
		paths = make([]Path, len(req.Path))
		query = Queries{
			Filter: make([]Filter_Query, len(req.Query.Filter)),
		}
		pages = int32(0)
	)

	for i := 0; i < len(req.Path); i++ {
		paths[i].CollectionID = req.Path[i].CollectionId
		paths[i].DocumentID = req.Path[i].DocumentId
		paths[i].NewDocument = req.Path[i].NewDocument
	}

	sorts := make([]Sort_Query, len(req.Query.Sort))
	for i := 0; i < len(req.Query.Sort); i++ {
		sorts[i] = Sort_Query{
			OrderBy:   req.Query.Sort[i].OrderBy,
			OrderType: OrderDir(req.Query.Sort[i].OrderType),
		}
	}

	if req.Query.DateRange != nil {
		query.DateRange = DateRange_Query{
			Field: req.Query.DateRange.Field,
			Start: req.Query.DateRange.Start.AsTime(),
			End:   req.Query.DateRange.End.AsTime(),
		}
	}

	for i := 0; i < len(req.Query.Filter); i++ {
		query.Filter[i] = Filter_Query{
			By: req.Query.Filter[i].By,
			Op: req.Query.Filter[i].Op,
		}

		if err = json.Unmarshal(req.Query.Filter[i].Val, &query.Filter[i].Val); err != nil {
			// TODO: handle error
			return nil, Queries{}, 0, err
		}
	}

	if req.IsPagination {
		pages = req.Pagination.Page
	}

	return paths, query, pages, err
}

func (srv *server) Retrive(ctx context.Context, req *collection_core_service.RetriveRequest) (*collection_core_service.RetriveResponse, error) {
	if req.Payload == nil {
		return &collection_core_service.RetriveResponse{
			Api: &collection_core_service.StandardAPIProto{
				Status:  "ERROR",
				Entity:  "retriveFirestoreDocument",
				State:   "retriveFirestoreDocumentError",
				Message: "Retrive Firestore Document Failed Read Source Data",
			},
		}, nil
	}

	paths, query, page, _ := srv.payloadBuilder(req.Payload)
	var (
		res = collection_core_service.RetriveResponse{}
		p   = Payload{
			RootCollection: req.Payload.RootCollection,
			RootDocument:   req.Payload.RootDocument,
			limit:          req.Payload.Limit,
			isPagination:   req.Payload.IsPagination,
			isDelete:       req.Payload.IsDelete,
			Data:           req.Payload.Data.AsMap(),
			Path:           paths,
			page:           page,
			query:          query,
		}
	)

	response, err := srv.source.Retrive(ctx, &p)
	if err != nil {
		res.Api = &collection_core_service.StandardAPIProto{
			Status:  "ERROR",
			Entity:  "retriveFirestoreDocument",
			State:   "retriveFirestoreDocumentError",
			Message: "Retrive Firestore Document Failed Read Source Data",
		}
		return &res, nil
	}

	data, err := json.Marshal(response.Data)
	if err != nil {
		res.Api = &collection_core_service.StandardAPIProto{
			Status:  "ERROR",
			Entity:  "retriveFirestoreDocument",
			State:   "retriveFirestoreDocumentMarshalResponseError",
			Message: "Retrive Firestore Document Failed Build Result Data",
		}
		return &res, nil
	}

	res.Api = &collection_core_service.StandardAPIProto{
		Status:  "SUCCESS",
		Entity:  "retriveFirestoreDocument",
		State:   "retriveFirestoreDocumentSuccess",
		Message: "Retrive Firestore Document Success",
		Meta: &collection_core_service.MetaProto{
			Page:    p.page,
			PerPage: p.limit,
			Size:    int32(response.Size),
		},
	}
	res.Data = data

	return &res, nil
}
