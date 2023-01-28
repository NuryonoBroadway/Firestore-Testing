package collectioncallgrpc

import (
	collectionxserver "firebaseapi/collectionx/collectionx_server"
	"testing"
)

func Test_Call_Grpc(t *testing.T) {
	config := collectionxserver.ServerConfig{
		ProjectID:            "privyfellowship-6a4e1",
		CredentialsFile:      "/home/slvr/FirebaseTestingApi/config/privyfellowship-6a4e1-firebase-adminsdk-g0c2d-66e133ed37.json",
		ProjectServiceConfig: collectionxserver.ProjectServiceConfig{},
	}

	CallGrpc(&config)
}
