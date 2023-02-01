package collectionxclient

import (
	"context"
	"encoding/json"
	"firebaseapi/collectionx/collection_core_service"
	"firebaseapi/helper"
	"fmt"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type Client interface {
	OpenConnection() (*client, error)
	Close() error
}

type client struct {
	cfg  *ClientConfig
	ctx  context.Context
	conn *grpc.ClientConn
	Documenter
}

func NewCollectionClient(cfg *ClientConfig) *client {
	return &client{cfg: cfg}
}

// In this method need context
func (c *client) NewProcess(ctx context.Context) (*client, error) {
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
		return nil, err
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

	pathProto, query, pagination := payloadBuilder(p)
	payloadProto := collection_core_service.PayloadProto{
		RootCollection: p.RootCollection,
		RootDocument:   p.RootDocument,
		Limit:          p.limit,
		IsPagination:   p.isPagination,
		Data:           structData,
		Path:           pathProto,
		Pagination:     pagination,
		Query:          query,
	}

	req := &collection_core_service.RetriveRequest{
		Payload: &payloadProto,
	}

	res, err := collection_core_service.NewServiceCollectionClient(p.conn).Retrive(p.ctx, req)
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
			Page:    res.Api.Meta.Page,
			PerPage: res.Api.Meta.PerPage,
			Total:   res.Api.Meta.Total,
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
