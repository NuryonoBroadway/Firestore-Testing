package collectionxserver

import (
	"context"
	"encoding/json"
	"errors"
	collectionxservice "firebaseapi/collectionx/collectionx_service"

	"cloud.google.com/go/firestore"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (srv *server) payloadBuilder(req *collectionxservice.PayloadProto) (path []Path, query Filtering, page int32) {
	path = make([]Path, len(req.Path))
	query = Filtering{
		Sort:   make([]Sort_Query, len(req.Query.Sort)),
		Filter: make([]Filter_Query, len(req.Query.Filter)),
	}

	for i := 0; i < len(req.Path); i++ {
		path[i].CollectionID = req.Path[i].CollectionId
		path[i].DocumentID = req.Path[i].DocumentId
		path[i].NewDocument = req.Path[i].NewDocument
		path[i].CollectionGroup = req.Path[i].CollectionGroup
	}

	for i := 0; i < len(req.Query.Sort); i++ {
		query.Sort[i] = Sort_Query{
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
		if err := json.Unmarshal(req.Query.Filter[i].Val, &query.Filter[i].Val); err != nil {
			return nil, Filtering{}, 0
		}
	}

	if req.IsPagination {
		page = req.Page
	}

	return path, query, page
}

func (srv *server) Retrive(ctx context.Context, req *collectionxservice.RetriveRequest) (*collectionxservice.RetriveResponse, error) {
	if req.Payload == nil {
		return nil, status.Error(codes.Unavailable, "payload unavailable")
	}

	paths, query, page := srv.payloadBuilder(req.Payload)
	var (
		res = collectionxservice.RetriveResponse{}
		p   = Payload{
			RootCollection: req.Payload.RootCollection,
			RootDocument:   req.Payload.RootDocument,
			limit:          req.Payload.Limit,
			isPagination:   req.Payload.IsPagination,
			isDelete:       req.Payload.IsDelete,
			Path:           paths,
			page:           page,
			query:          query,
		}
	)

	retriveRes, err := srv.source.Retrive(ctx, &p)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	data, err := json.Marshal(retriveRes.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res.Api = &collectionxservice.StandardAPIProto{
		Message: "Retrive Firestore Document Success",
		Meta: &collectionxservice.MetaProto{
			Page:    p.page,
			PerPage: p.limit,
			Total:   int32(retriveRes.Size),
		},
	}
	res.Data = data

	return &res, nil
}

func (srv *server) Snapshots(req *collectionxservice.SnapshotsRequest, stream collectionxservice.ServiceCollection_SnapshotsServer) error {
	if req.Payload == nil {
		return status.Error(codes.Unavailable, "payload unavailable")
	}

	paths, query, page := srv.payloadBuilder(req.Payload)
	var (
		p = Payload{
			RootCollection: req.Payload.RootCollection,
			RootDocument:   req.Payload.RootDocument,
			limit:          req.Payload.Limit,
			isPagination:   req.Payload.IsPagination,
			isDelete:       req.Payload.IsDelete,
			Path:           paths,
			page:           page,
			query:          query,
		}
	)

	col, doc, err := srv.source.Snapshots(stream.Context(), &p)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if col != nil {
		defer col.Stop()

		for {
			snap, err := col.Next()
			if err != nil {
				switch status.Code(err) {
				case codes.Canceled:
					return status.Error(codes.Canceled, "context canceled")
				case codes.DeadlineExceeded:
					return status.Error(codes.DeadlineExceeded, "context deadline exceeded")
				default:
					return status.Error(codes.Internal, err.Error())
				}
			}

			var (
				response = &collectionxservice.SnapshotsResponse{
					Api: &collectionxservice.StandardAPIProto{
						Message: "Snapshots Firestore Collection Success",
					},
				}
			)

			if snap != nil {
				for _, change := range snap.Changes {
					data, err := json.Marshal(change.Doc.Data())
					if err != nil {
						return status.Error(codes.Internal, err.Error())
					}

					response.DocumentChange = &collectionxservice.DocumentChange{
						Data: data,
						Timestamp: &collectionxservice.TimestampProto{
							CreatedTime: timestamppb.New(change.Doc.CreateTime),
							ReadTime:    timestamppb.New(change.Doc.ReadTime),
							UpdateTime:  timestamppb.New(change.Doc.UpdateTime),
						},
					}

					switch change.Kind {
					case firestore.DocumentAdded:
						response.DocumentChange.Kind = collectionxservice.DocumentChangeKind_DOCUMENT_KIND_ADDED
						if err := stream.Send(response); err != nil {
							return err
						}

					case firestore.DocumentModified:
						response.DocumentChange.Kind = collectionxservice.DocumentChangeKind_DOCUMENT_KIND_MODIFIED
						if err := stream.Send(response); err != nil {
							return err
						}
					case firestore.DocumentRemoved:
						response.DocumentChange.Kind = collectionxservice.DocumentChangeKind_DOCUMENT_KIND_REMOVED
						if err := stream.Send(response); err != nil {
							return err
						}
					}
				}
			}
		}
	} else if doc != nil {
		defer doc.Stop()

		for {
			snap, err := doc.Next()
			if err != nil {
				if err != nil {
					switch status.Code(err) {
					case codes.Canceled:
						return status.Error(codes.Canceled, "context canceled")
					case codes.DeadlineExceeded:
						return status.Error(codes.DeadlineExceeded, "context deadline exceeded")
					default:
						return status.Error(codes.Internal, err.Error())
					}
				}
			}

			var (
				response = &collectionxservice.SnapshotsResponse{
					Api: &collectionxservice.StandardAPIProto{
						Message: "Snapshots Firestore Document Success",
					},
				}
			)

			if !snap.Exists() {
				return status.Error(codes.NotFound, "document not found")
			}

			data, err := json.Marshal(snap.Data())
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			response.DocumentChange = &collectionxservice.DocumentChange{
				Kind: collectionxservice.DocumentChangeKind_DOCUMENT_KIND_SNAPSHOTS,
				Data: data,
				Timestamp: &collectionxservice.TimestampProto{
					CreatedTime: timestamppb.New(snap.CreateTime),
					ReadTime:    timestamppb.New(snap.ReadTime),
					UpdateTime:  timestamppb.New(snap.UpdateTime),
				},
			}

			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}

	return errors.New("not-implemented")
}
