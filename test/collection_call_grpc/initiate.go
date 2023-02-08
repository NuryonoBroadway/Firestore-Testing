package collectioncallgrpc

import (
	"context"
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	collectionxserver "firebaseapi/collectionx/collectionx_server"
	"firebaseapi/helper"
	"net"

	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func CallGrpc(cfg *collectionxserver.ServerConfig) *grpc.Server {
	ctx := context.Background()

	srv := collectionxserver.NewServer(
		ctx, 
		cfg, 
		collectionxserver.WithMaxConcurrent(5),
		collectionxserver.WithSubscribeAsync(true),
		collectionxserver.WithSubscriberID("PrivyFlowSe"),
	)

	defer srv.GracefulStop()

	logger.Infof("starting privypass-collection-core-se grpc services... 0.0.0.0:9090")

	listen, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		logger.Warnf("cannot listen grpc port, err: %v", err.Error())
	}

	if err := srv.Serve(listen); err != nil {
		logger.Fatalf("service  grpc stopped, err: %v", err.Error())
	}

	return srv
}

type Sort struct {
	By  string                     `json:"by"`
	Dir collectionxclient.OrderDir `json:"dir"`
}

type Filter struct {
	By  string          `json:"by"`
	Op  helper.Operator `json:"op"`
	Val interface{}     `json:"val"`
}
