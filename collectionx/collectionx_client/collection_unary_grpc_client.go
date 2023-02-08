package collectionxclient

import (
	"context"
	"encoding/json"
	"errors"
	collectionxservice "firebaseapi/collectionx/collectionx_service"
	"firebaseapi/helper"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type Client interface {
	OpenConnection(ctx context.Context) (*client, error)
	Close() error
}

type client struct {
	cfg   *configClient
	ctx   context.Context
	conn  *grpc.ClientConn
	topic *pubsub.Topic

	logf *logrus.Logger

	payload *Payload
}

func NewCollectionClient(cfg *configClient) Client {
	return &client{cfg: cfg, logf: logrus.New()}
}

func (c *client) OpenConnection(ctx context.Context) (*client, error) {

	if c.cfg == nil {
		return nil, fmt.Errorf("config is empty")
	}

	pubsub, err := pubsub.NewClient(ctx, c.cfg.ProjectName, option.WithCredentialsFile(c.cfg.pubsubCredential))
	if err != nil {
		return nil, fmt.Errorf("config is empty")
	}

	topic := pubsub.Topic(c.cfg.PubSubTopic)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("topic not registered")
	}

	connection, err := grpc.DialContext(
		ctx, c.cfg.GrpcAddress,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithKeepaliveParams(kacp),
	)
	if err != nil {
		return nil, err
	}

	c = &client{
		topic: topic,
		conn:  connection,
		ctx:   ctx,
	}
	c.payload = &Payload{
		client: c,
	}

	return c, nil
}

func (c *client) Close() error {
	if c.topic != nil && c.conn != nil {
		c.topic.Stop()
		return c.conn.Close()
	}

	return nil
}

func (c *client) Col(path string) Collector {
	c.payload.Path = append(c.payload.Path, Path{
		CollectionID: path,
	})

	return c.payload
}

func (c *client) Doc(path string) Documenter {
	c.payload.Path = append(c.payload.Path, Path{
		DocumentID: path,
	})

	return c.payload
}

func (c *client) ColGroup(id string) CollectorGroup {
	c.payload.Path = append(c.payload.Path, Path{
		CollectionGroup: id,
	})

	return c.payload
}

func (p *Payload) Retrive() (response *StandardAPI, err error) {
	response = new(StandardAPI)
	structData, err := structpb.NewStruct(map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed build collection core request data: %w", err)
	}

	path, query, err := payloadBuilder(p)
	if err != nil {
		return nil, err
	}

	payloadProto := collectionxservice.PayloadProto{
		RootCollection: p.RootCollection,
		RootDocument:   p.RootDocument,
		Limit:          p.limit,
		IsPagination:   p.isPagination,
		Data:           structData,
		Path:           path,
		Page:           p.page,
		Query:          query,
	}

	req := &collectionxservice.RetriveRequest{
		Payload: &payloadProto,
	}

	res, err := collectionxservice.NewServiceCollectionClient(p.client.conn).Retrive(p.client.ctx, req, grpc.WaitForReady(true))
	if err != nil {
		response.StandardAPIDefault = NewErrorResponse()
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.DeadlineExceeded:
				response.StandardAPIDefault = response.WithRepresentative(helper.DEADLINE_EXCEEDED, status.Message())
			case codes.Canceled:
				response.StandardAPIDefault = response.WithRepresentative(helper.CANCELLED, status.Message())
			case codes.Internal:
				response.StandardAPIDefault = response.WithRepresentative(helper.INTERNAL, status.Message())
			case codes.NotFound:
				response.StandardAPIDefault = response.WithRepresentative(helper.NOT_FOUND, status.Message())
			default:
				response.StandardAPIDefault = response.WithRepresentative(helper.FAILED_PRECONDITION, fmt.Errorf("failed retrive collection core: %w", err).Error())
			}

			return response, err
		}
	}

	response.StandardAPIDefault = NewSuccessResponse().WithRepresentative(helper.SUCCESS, res.Api.Message)
	if p.Path[len(p.Path)-1].CollectionID != "" || p.Path[len(p.Path)-1].CollectionGroup != "" {
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
		Page:    res.Api.Meta.Page,
		PerPage: res.Api.Meta.PerPage,
		Total:   res.Api.Meta.Total,
	}

	return response, nil
}

func (p *Payload) Save() (response *StandardAPIDefault, err error) {
	if len(p.Data.Row) < 1 {
		return NewErrorResponse().
			WithRepresentative(helper.UNAVAILABLE, "data len is 0"), errors.New("len data is 0")
	}

	data, err := json.Marshal(&p)
	if err != nil {
		return NewErrorResponse().
			WithRepresentative(helper.INTERNAL, "cannot marshal object"), err
	}

	id, err := p.client.topic.Publish(p.client.ctx, &pubsub.Message{
		Data: data,
	}).Get(p.client.ctx)
	if err != nil {
		return NewErrorResponse().
			WithRepresentative(helper.UNAVAILABLE, "failed to get response"), err
	}

	return NewSuccessResponse().
		WithRepresentative(helper.SUCCESS, fmt.Sprintf("success publish message: %v", id)), nil
}

func (s *CollectionxSnapshots) message() bool {
	res, err := s.ws.Recv()
	if err != nil {
		if err == io.EOF {
			s.err = io.EOF
			return false
		}

		s.err = err
		return false
	}

	if res != nil {
		s.res = res
	}
	return false
}

func (s *CollectionxSnapshots) Receive() (snap *Snapshots, err error) {
	snap = new(Snapshots)
	if s.err != nil {
		return nil, s.err
	}

	for s.message() {
	}

	if s.err != nil {
		s.Close()
		if s.err == io.EOF {
			return nil, s.err
		}

		erros := status.Convert(s.err)
		snap.StandardAPIDefault = NewErrorResponse()

		switch erros.Code() {
		case codes.Canceled:
			snap.StandardAPIDefault = snap.WithRepresentative(helper.CANCELLED, erros.Message())
		case codes.Internal:
			snap.StandardAPIDefault = snap.WithRepresentative(helper.INTERNAL, erros.Message())
		case codes.DeadlineExceeded:
			snap.StandardAPIDefault = snap.WithRepresentative(helper.DEADLINE_EXCEEDED, erros.Message())
		case codes.NotFound:
			snap.StandardAPIDefault = snap.WithRepresentative(helper.NOT_FOUND, erros.Message())
		default:
			snap.StandardAPIDefault = snap.WithRepresentative(helper.FAILED_PRECONDITION, fmt.Errorf("failed retrive collection core: %w", err).Error())
		}

		return snap, s.err
	}

	snap.StandardAPIDefault = NewSuccessResponse().WithRepresentative(helper.SUCCESS, s.res.Api.Message)
	snap.Timestamp = Timestamp{
		CreatedTime: s.res.DocumentChange.Timestamp.CreatedTime.AsTime(),
		ReadTime:    s.res.DocumentChange.Timestamp.ReadTime.AsTime(),
		UpdateTime:  s.res.DocumentChange.Timestamp.UpdateTime.AsTime(),
	}

	d := make(map[string]interface{})
	if err := json.Unmarshal(s.res.DocumentChange.Data, &d); err != nil {
		return nil, err
	}

	if s.isCol {
		snap.Data.Type = helper.Collection
	} else {
		snap.Data.Type = helper.Document
	}

	if d != nil {
		snap.Data.Data = d
	}

	switch s.res.DocumentChange.Kind {
	case collectionxservice.DocumentChangeKind_DOCUMENT_KIND_ADDED:
		snap.Kind = DOCUMENT_KIND_ADDED.ToString()
	case collectionxservice.DocumentChangeKind_DOCUMENT_KIND_REMOVED:
		snap.Kind = DOCUMENT_KIND_REMOVED.ToString()
	case collectionxservice.DocumentChangeKind_DOCUMENT_KIND_MODIFIED:
		snap.Kind = DOCUMENT_KIND_MODIFIED.ToString()
	case collectionxservice.DocumentChangeKind_DOCUMENT_KIND_SNAPSHOTS:
		snap.Kind = DOCUMENT_KIND_SNAPSHOTS.ToString()
	}

	return snap, nil
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
	pathProto, query, err := payloadBuilder(p)
	if err != nil {
		return nil, err
	}

	payloadProto := collectionxservice.PayloadProto{
		RootCollection: p.RootCollection,
		RootDocument:   p.RootDocument,
		Limit:          p.limit,
		IsPagination:   p.isPagination,
		Path:           pathProto,
		Page:           p.page,
		Query:          query,
	}

	req := &collectionxservice.SnapshotsRequest{
		Payload: &payloadProto,
	}

	stream, err := collectionxservice.NewServiceCollectionClient(p.client.conn).Snapshots(p.client.ctx, req, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	return &CollectionxSnapshots{
		isCol:  p.Path[len(p.Path)-1].CollectionID != "",
		ws:     stream,
		client: p.client,
	}, nil
}
