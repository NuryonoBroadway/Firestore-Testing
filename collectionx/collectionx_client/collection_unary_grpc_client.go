package collectionxclient

import (
	"context"
	"encoding/json"
	collectionxservice "firebaseapi/collectionx/collectionx_service"
	"firebaseapi/helper"
	"fmt"
	"io"
	"reflect"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client interface {
	OpenConnection() (*client, error)
	Close() error
}

type client struct {
	cfg    *configClient
	ctx    context.Context
	cancel context.CancelFunc
	conn   *grpc.ClientConn
	Documenter
}

func NewCollectionClient(cfg *configClient) Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &client{cfg: cfg, ctx: ctx, cancel: cancel}
}

func (c *client) OpenConnection() (*client, error) {
	if c.cfg == nil {
		return nil, fmt.Errorf("config is empty")
	}

	conn, err := grpc.DialContext(
		c.ctx, c.cfg.GrpcAddress,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		c.cancel()
		return nil, fmt.Errorf("error when connect ")
	}

	c.conn = conn
	c.Documenter = NewCollectionPayloads(
		WithRootCollection(c.cfg.ProjectRootCollection),
		WithRootDocuments(c.cfg.ProjectRootDocument),
		WithGRPCCon(conn),
		WithContext(c.ctx),
	)

	return c, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (p *Payload) Retrive() (*StandardAPI, error) {
	structData, err := structpb.NewStruct(p.Data)
	if err != nil {
		return nil, fmt.Errorf("failed build collection core request data: %w", err)
	}

	pathProto := make([]*collectionxservice.PathProto, len(p.Path))
	for i := 0; i < len(p.Path); i++ {
		pathProto[i] = &collectionxservice.PathProto{
			CollectionId: p.Path[i].CollectionID,
			DocumentId:   p.Path[i].DocumentID,
			NewDocument:  p.Path[i].NewDocument,
		}
	}

	filters := make([]*collectionxservice.FilterProto, len(p.query.Filter))
	if len(filters) != 0 {
		for i := 0; i < len(p.query.Filter); i++ {
			filters[i] = &collectionxservice.FilterProto{
				By: p.query.Filter[i].By,
				Op: p.query.Filter[i].Op,
			}

			xtyp := reflect.TypeOf(p.query.Filter[i].Val)
			switch xtyp.Kind() {
			case reflect.Bool:
				filters[i].Val = &collectionxservice.FilterProto_ValBool{
					ValBool: p.query.Filter[i].Val.(bool),
				}
			case reflect.String:
				filters[i].Val = &collectionxservice.FilterProto_ValString{
					ValString: p.query.Filter[i].Val.(string),
				}
			case reflect.Int64:
				filters[i].Val = &collectionxservice.FilterProto_ValInt{
					ValInt: p.query.Filter[i].Val.(int64),
				}

			}

		}
	}

	sorts := new(collectionxservice.SortProto)
	if p.query.Sort.OrderBy != "" {
		sorts.OrderBy = p.query.Sort.OrderBy
		sorts.OrderType = collectionxservice.OrderTypeProto(p.query.Sort.OrderType)
	} else {
		sorts.OrderBy = "created_at"
		sorts.OrderType = collectionxservice.OrderTypeProto(Asc)
	}

	ranges := new(collectionxservice.DateRangeProto)
	if p.query.DateRange.Field != "" {
		ranges.Field = p.query.DateRange.Field
		ranges.Start = timestamppb.New(p.query.DateRange.Start)
		ranges.End = timestamppb.New(p.query.DateRange.End)
	}

	query := &collectionxservice.FilteringProto{
		Sort:      sorts,
		Filter:    filters,
		DateRange: ranges,
	}

	pagination := new(collectionxservice.PaginationProto)
	if p.isPagination {
		if p.limit == 0 {
			p.limit = 2 // default limit
		}
		pagination.Page = p.pagination.Page
	}

	payloadProto := collectionxservice.PayloadProto{
		RootCollection: p.RootCollection,
		RootDocument:   p.RootDocument,
		Limit:          p.limit,
		IsPagination:   p.isPagination,
		IsDelete:       p.isDelete,
		Data:           structData,
		Path:           pathProto,
		Pagination:     pagination,
		Query:          query,
	}

	req := &collectionxservice.RetriveRequest{
		Payload: &payloadProto,
	}

	res, err := collectionxservice.NewServiceCollectionClient(p.conn).Retrive(p.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed retrive collection core: %w", err)
	}

	isCol := p.Path[len(p.Path)-1].CollectionID != ""
	resp := StandardAPIDefault{
		Status:  res.Api.Status,
		Entity:  res.Api.Entity,
		State:   res.Api.State,
		Message: res.Api.Message,
	}

	response := StandardAPI{
		StandardAPIDefault: resp,
	}

	if res.Api.Error == nil {
		if isCol {
			d := []ListValue{}
			if err := json.Unmarshal(res.Data, &d); err != nil {
				return nil, fmt.Errorf("failed json unmarshal payload data collection core: %w", err)
			}
			response.Data = Data{
				Type: helper.Collection,
				Data: d,
			}
		} else {
			d := make(map[string]interface{})
			if err := json.Unmarshal(res.Data, &d); err != nil {
				return nil, fmt.Errorf("failed json unmarshal payload data collection core: %w", err)
			}
			response.Data = Data{
				Type: helper.Document,
				Data: d,
			}
		}

		response.Meta = Meta{
			Page:      res.Api.Meta.Page,
			PerPage:   res.Api.Meta.PerPage,
			Total:     res.Api.Meta.Total,
			OrderBy:   res.Api.Meta.OrderBy,
			OrderType: res.Api.Meta.OrderType,
		}
	}

	if res.Api.Error != nil {
		buildErr := Error{}
		buildErr.General = res.Api.Error.General
		buildErr.Validation = make([]map[string]string, len(res.Api.Error.Validation))
		for i, v := range res.Api.Error.Validation {
			buildErr.Validation[i] = map[string]string{
				v.Key: v.Value,
			}
		}
		resp.Error = &buildErr
	}
	return &response, nil
}

func (s *CollectionxSnapshots) message() bool {
	res, err := s.ws.Recv()
	if err != nil {
		if err == io.EOF {
			s.err = io.EOF
			return false
		}

		s.err = fmt.Errorf("failed retrive collection core: %w", err)
		return false
	}

	resp := StandardAPIDefault{
		Status:  res.Api.Status,
		Entity:  res.Api.Entity,
		State:   res.Api.State,
		Message: res.Api.Message,
	}

	if res.Api.Error == nil {
		s.snapshots.StandardAPIDefault = resp
		s.snapshots.Timestamp = Timestamp{
			CreatedTime: res.DocumentChange.Timestamp.CreatedTime.AsTime(),
			ReadTime:    res.DocumentChange.Timestamp.ReadTime.AsTime(),
			UpdateTime:  res.DocumentChange.Timestamp.UpdateTime.AsTime(),
		}

		d := make(map[string]interface{})
		if err := json.Unmarshal(res.DocumentChange.Data, &d); err != nil {
			s.err = fmt.Errorf("failed json unmarshal payload data collection core: %w", err)
			return false
		}

		if s.isCol {
			s.snapshots.Data.Type = helper.Collection
		} else {
			s.snapshots.Data.Type = helper.Document
		}

		if d != nil {
			s.snapshots.Data.Data = d
		}

		switch res.DocumentChange.Kind {
		case collectionxservice.DocumentChangeKind_DOCUMENT_KIND_ADDED:
			s.snapshots.Kind = DOCUMENT_KIND_ADDED.ToString()
		case collectionxservice.DocumentChangeKind_DOCUMENT_KIND_REMOVED:
			s.snapshots.Kind = DOCUMENT_KIND_REMOVED.ToString()
		case collectionxservice.DocumentChangeKind_DOCUMENT_KIND_MODIFIED:
			s.snapshots.Kind = DOCUMENT_KIND_MODIFIED.ToString()
		case collectionxservice.DocumentChangeKind_DOCUMENT_KIND_SNAPSHOTS:
			s.snapshots.Kind = DOCUMENT_KIND_SNAPSHOTS.ToString()
		}
	}

	if res.Api.Error != nil {
		buildErr := Error{}
		buildErr.General = res.Api.Error.General
		buildErr.Validation = make([]map[string]string, len(res.Api.Error.Validation))
		for i, v := range res.Api.Error.Validation {
			buildErr.Validation[i] = map[string]string{
				v.Key: v.Value,
			}
		}
		resp.Error = &buildErr
		s.snapshots.StandardAPIDefault = resp
	}

	return false
}

func (s *CollectionxSnapshots) Receive() (*Snapshots, error) {
	if s.err != nil {
		return nil, s.err
	}

	for s.message() {
	}

	if s.err == io.EOF {
		return nil, s.err
	}

	if s.err != nil {
		s.Close()
		return nil, s.err
	}

	return &Snapshots{
		StandardAPIDefault: s.snapshots.StandardAPIDefault,
		Kind:               s.snapshots.Kind,
		Data:               s.snapshots.Data,
		Timestamp:          s.snapshots.Timestamp,
	}, nil
}

func (s *CollectionxSnapshots) Close() {
	err := s.ws.CloseSend()
	if s.err != nil { // don't change existing error
		return
	}
	if err != nil {
		// if an error occurs while closing the stream
		s.err = err
		return
	}
}

func (p *Payload) Snapshots() (*CollectionxSnapshots, error) {
	structData, err := structpb.NewStruct(p.Data)
	if err != nil {
		return nil, fmt.Errorf("failed build collection core request data: %w", err)
	}

	pathProto := make([]*collectionxservice.PathProto, len(p.Path))
	for i := 0; i < len(p.Path); i++ {
		pathProto[i] = &collectionxservice.PathProto{
			CollectionId: p.Path[i].CollectionID,
			DocumentId:   p.Path[i].DocumentID,
			NewDocument:  p.Path[i].NewDocument,
		}
	}

	filters := make([]*collectionxservice.FilterProto, len(p.query.Filter))
	if len(filters) != 0 {
		for i := 0; i < len(p.query.Filter); i++ {
			filters[i] = &collectionxservice.FilterProto{
				By: p.query.Filter[i].By,
				Op: p.query.Filter[i].Op,
			}

			xtyp := reflect.TypeOf(p.query.Filter[i].Val)
			switch xtyp.Kind() {
			case reflect.Bool:
				filters[i].Val = &collectionxservice.FilterProto_ValBool{
					ValBool: p.query.Filter[i].Val.(bool),
				}
			case reflect.String:
				filters[i].Val = &collectionxservice.FilterProto_ValString{
					ValString: p.query.Filter[i].Val.(string),
				}
			case reflect.Int64:
				filters[i].Val = &collectionxservice.FilterProto_ValInt{
					ValInt: p.query.Filter[i].Val.(int64),
				}

			}

		}
	}

	sorts := new(collectionxservice.SortProto)
	if p.query.Sort.OrderBy != "" {
		sorts.OrderBy = p.query.Sort.OrderBy
		sorts.OrderType = collectionxservice.OrderTypeProto(p.query.Sort.OrderType)
	} else {
		sorts.OrderBy = "created_at"
		sorts.OrderType = collectionxservice.OrderTypeProto(Asc)
	}

	ranges := new(collectionxservice.DateRangeProto)
	if p.query.DateRange.Field != "" {
		ranges.Field = p.query.DateRange.Field
		ranges.Start = timestamppb.New(p.query.DateRange.Start)
		ranges.End = timestamppb.New(p.query.DateRange.End)
	}

	query := &collectionxservice.FilteringProto{
		Sort:      sorts,
		Filter:    filters,
		DateRange: ranges,
	}

	pagination := new(collectionxservice.PaginationProto)
	if p.isPagination {
		if p.limit == 0 {
			p.limit = 2 // default limit
		}
		pagination.Page = p.pagination.Page
	}

	payloadProto := collectionxservice.PayloadProto{
		RootCollection: p.RootCollection,
		RootDocument:   p.RootDocument,
		Limit:          p.limit,
		IsPagination:   p.isPagination,
		IsDelete:       p.isDelete,
		Data:           structData,
		Path:           pathProto,
		Pagination:     pagination,
		Query:          query,
	}

	req := &collectionxservice.SnapshotsRequest{
		Payload: &payloadProto,
	}

	stream, err := collectionxservice.NewServiceCollectionClient(p.conn).Snapshots(p.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed retrive stream collection core: %w", err)
	}

	return &CollectionxSnapshots{
		isCol: p.Path[len(p.Path)-1].CollectionID != "",
		ws:    stream,
	}, nil
}
