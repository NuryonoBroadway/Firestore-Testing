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

	if req.Payload.Query.Sort != nil {
		query.Sort = Sort_Query{
			OrderBy:   req.Payload.Query.Sort.OrderBy,
			OrderType: OrderDir(req.Payload.Query.Sort.OrderType.Number()),
		}
	}

	if req.Payload.Query.DateRange != nil {
		query.DateRange = DateRange_Query{
			Field: req.Payload.Query.DateRange.Field,
			Start: req.GetPayload().Query.DateRange.Start.AsTime(),
			End:   req.GetPayload().Query.DateRange.End.AsTime(),
		}
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

	if req.Payload.IsPagination {
		pagination = Pagination{
			Page: req.Payload.Pagination.Page,
		}
	}

	var (
		res = collectionxservice.RetriveResponse{}
		p   = Payload{
			RootCollection: req.Payload.RootCollection,
			RootDocument:   req.Payload.RootDocument,
			limit:          req.Payload.Limit,
			isPagination:   req.Payload.IsPagination,
			isDelete:       req.Payload.IsDelete,
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

	data, err := json.Marshal(retriveRes.Data)
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
			Total:     int32(retriveRes.Total),
			OrderBy:   p.query.Sort.OrderBy,
			OrderType: p.query.Sort.OrderType.ToString(),
		},
	}
	res.Data = data

	return &res, nil
}

func (srv *server) Snapshots(req *collectionxservice.SnapshotsRequest, stream collectionxservice.ServiceCollection_SnapshotsServer) error {
	if req.Payload == nil {
		return stream.Send(&collectionxservice.SnapshotsResponse{
			Api: &collectionxservice.StandardAPIProto{
				Status:  "ERROR",
				Entity:  "snapshotsFirestore",
				State:   "snapshotsirestoreError",
				Message: "Snapshots Firestore Failed Read Source Data",
			},
		})
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

	if req.Payload.Query.Sort != nil {
		query.Sort = Sort_Query{
			OrderBy:   req.Payload.Query.Sort.OrderBy,
			OrderType: OrderDir(req.Payload.Query.Sort.OrderType.Number()),
		}
	}

	if req.Payload.Query.DateRange != nil {
		query.DateRange = DateRange_Query{
			Field: req.Payload.Query.DateRange.Field,
			Start: req.GetPayload().Query.DateRange.Start.AsTime(),
			End:   req.GetPayload().Query.DateRange.End.AsTime(),
		}
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

	if req.Payload.IsPagination {
		pagination = Pagination{
			Page: req.Payload.Pagination.Page,
		}
	}

	var (
		p = Payload{
			RootCollection: req.Payload.RootCollection,
			RootDocument:   req.Payload.RootDocument,
			limit:          req.Payload.Limit,
			isPagination:   req.Payload.IsPagination,
			isDelete:       req.Payload.IsDelete,
			Data:           req.Payload.Data.AsMap(),
			Path:           paths,
			pagination:     pagination,
			query:          query,
		}
	)

	col, doc, err := srv.source.Snapshots(stream.Context(), &p)
	if err != nil {
		return stream.Send(&collectionxservice.SnapshotsResponse{
			Api: &collectionxservice.StandardAPIProto{
				Status:  "ERROR",
				Entity:  "snapshotsFirestore",
				State:   "snapshotsFirestoreError",
				Message: "Snapshots Firestore Failed Start Snapshots",
				Error: &collectionxservice.ErrorProto{
					General: err.Error(),
				},
			},
		})
	}

	if col != nil {
		defer col.Stop()

		for {
			snap, err := col.Next()
			if err != nil {
				if e := status.Code(err); e == codes.Canceled || e == codes.DeadlineExceeded {
					return stream.Send(&collectionxservice.SnapshotsResponse{
						Api: &collectionxservice.StandardAPIProto{
							Status:  "ERROR",
							Entity:  "snapshotsFirestoreCollection",
							State:   "snapshotsFirestoreCollectionError",
							Message: "Snapshots Firestore Collection Failed Context Error",
							Error: &collectionxservice.ErrorProto{
								General: err.Error(),
							},
						},
					})
				}

				return stream.Send(&collectionxservice.SnapshotsResponse{
					Api: &collectionxservice.StandardAPIProto{
						Status:  "ERROR",
						Entity:  "snapshotsFirestoreCollection",
						State:   "snapshotsFirestoreCollectionError",
						Message: "Snapshots Firestore Collection Failed Error Found",
						Error: &collectionxservice.ErrorProto{
							General: err.Error(),
						},
					},
				})
			}

			var (
				response = &collectionxservice.SnapshotsResponse{
					Api: &collectionxservice.StandardAPIProto{
						Status:  "SUCCESS",
						Entity:  "snapshotsFirestoreCollection",
						State:   "snapshotsFirestoreCollectionSuccess",
						Message: "Snapshots Firestore Collection Success",
					},
				}
			)

			if snap != nil {
				for _, change := range snap.Changes {
					data, err := json.Marshal(change.Doc.Data())
					if err != nil {
						return stream.Send(&collectionxservice.SnapshotsResponse{
							Api: &collectionxservice.StandardAPIProto{
								Status:  "ERROR",
								Entity:  "retriveFirestoreCollection",
								State:   "retriveFirestoreCollectionMarshalResponseError",
								Message: "Retrive Firestore Collection Failed Build Result Data",
								Error: &collectionxservice.ErrorProto{
									General: err.Error(),
								},
							}})
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
				if e := status.Code(err); e == codes.Canceled || e == codes.DeadlineExceeded {
					return stream.Send(&collectionxservice.SnapshotsResponse{
						Api: &collectionxservice.StandardAPIProto{
							Status:  "ERROR",
							Entity:  "snapshotsFirestoreDocument",
							State:   "snapshotsFirestoreDocumentError",
							Message: "Snapshots Firestore Document Failed Context Error",
							Error: &collectionxservice.ErrorProto{
								General: err.Error(),
							},
						},
					})
				}

				return stream.Send(&collectionxservice.SnapshotsResponse{
					Api: &collectionxservice.StandardAPIProto{
						Status:  "ERROR",
						Entity:  "snapshotsFirestoreDocument",
						State:   "snapshotsFirestoreDocumentError",
						Message: "Snapshots Firestore Document Failed Error Found",
						Error: &collectionxservice.ErrorProto{
							General: err.Error(),
						},
					},
				})
			}

			var (
				response = &collectionxservice.SnapshotsResponse{
					Api: &collectionxservice.StandardAPIProto{
						Status:  "SUCCESS",
						Entity:  "snapshotsFirestoreDocument",
						State:   "snapshotsFirestoreDocumentSuccess",
						Message: "Snapshots Firestore Document Success",
					},
				}
			)

			if !snap.Exists() {
				return stream.Send(&collectionxservice.SnapshotsResponse{
					Api: &collectionxservice.StandardAPIProto{
						Status:  "ERROR",
						Entity:  "snapshotsFirestoreDocument",
						State:   "snapshotsFirestoreDocumentError",
						Message: "Snapshots Firestore Document Failed Error Found",
						Error: &collectionxservice.ErrorProto{
							General: errors.New("document no longer exists").Error(),
						},
					},
				})
			}

			data, err := json.Marshal(snap.Data())
			if err != nil {
				return stream.Send(&collectionxservice.SnapshotsResponse{
					Api: &collectionxservice.StandardAPIProto{
						Status:  "ERROR",
						Entity:  "snapshotsFirestoreDocument",
						State:   "snapshotsFirestoreDocumentMarshalResponseError",
						Message: "Retrive Firestore Document Failed Build Result Data",
						Error: &collectionxservice.ErrorProto{
							General: err.Error(),
						},
					}})
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
