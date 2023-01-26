package collectionxclient

import (
	"context"
	"encoding/json"
	collectionxservice "firebaseapi/collectionx/collectionx_service"
	"firebaseapi/helper"
	"fmt"
	reflect "reflect"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	structpb "google.golang.org/protobuf/types/known/structpb"
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
		WithRootCollection(string(c.cfg.ProjectRootCollection)),
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

	filters := make([]*collectionxservice.FilterProto, len(p.filter))
	for i := 0; i < len(p.filter); i++ {
		filters[i] = &collectionxservice.FilterProto{
			By: p.filter[i].By,
			Op: p.filter[i].Op,
		}

		xtyp := reflect.TypeOf(p.filter[i].Val)
		switch xtyp.Kind() {
		case reflect.Bool:
			filters[i].Val = &collectionxservice.FilterProto_ValBool{
				ValBool: p.filter[i].Val.(bool),
			}
		case reflect.String:
			filters[i].Val = &collectionxservice.FilterProto_ValString{
				ValString: p.filter[i].Val.(string),
			}
		case reflect.Int64:
			filters[i].Val = &collectionxservice.FilterProto_ValInt{
				ValInt: p.filter[i].Val.(int64),
			}

		}
	}

	sorts := &collectionxservice.SortProto{
		By:  p.sort.By,
		Dir: p.sort.Dir,
	}

	payloadProto := collectionxservice.PayloadProto{
		RootCollection: p.RootCollection,
		Filter:         filters,
		Limit:          p.limit,
		Sort:           sorts,
		IsDelete:       p.IsDelete,
		Data:           structData,
		Path:           pathProto,
	}

	req := &collectionxservice.RetriveRequest{
		Payload: &payloadProto,
	}

	res, err := collectionxservice.NewServiceCollectionClient(p.conn).Retrive(p.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed retrive collection core: %w", err)
	}

	isCol := p.Path[len(p.Path)-1].CollectionID != ""
	resp := StandardAPI{
		Status:  res.Api.Status,
		Entity:  res.Api.Entity,
		State:   res.Api.State,
		Message: res.Api.Message,
	}
	if isCol {
		d := []ListValue{}
		if err := json.Unmarshal(res.Data, &d); err != nil {
			return nil, fmt.Errorf("failed json unmarshal payload data collection core: %w", err)
		}
		resp.Data = Data{
			Type: helper.Collection,
			Data: d,
		}
	} else {
		d := make(map[string]interface{})
		if err := json.Unmarshal(res.Data, &d); err != nil {
			return nil, fmt.Errorf("failed json unmarshal payload data collection core: %w", err)
		}
		resp.Data = Data{
			Type: helper.Document,
			Data: d,
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
	return &resp, nil
}
