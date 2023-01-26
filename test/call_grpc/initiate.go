package callgrpc

import (
	"firebaseapi/collectionx"
	"net"

	logger "github.com/sirupsen/logrus"
)

func CallGrpc(cfg *collectionx.Config) {
	collx := collectionx.NewCollectionCore_SourceDocument(cfg)
	srv := collectionx.NewServer(collx)
	defer srv.GracefulStop()

	logger.Infof("starting privypass-collection-core-se grpc services... %v", cfg.GrpcAddress)

	listen, err := net.Listen("tcp", cfg.GrpcAddress)
	if err != nil {
		logger.Warnf("cannot listen grpc port, err: %v", err.Error())
	}

	if err := srv.Serve(listen); err != nil {
		logger.Fatalf("service  grpc stopped, err: %v", err.Error())
	}
}
