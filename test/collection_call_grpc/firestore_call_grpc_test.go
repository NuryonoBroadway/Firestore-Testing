package collectioncallgrpc

import (
	collection_core_server "firebaseapi/collectionx/collectionx_server"
	"testing"
)

func Test_Call_Grpc(t *testing.T) {
	config := collection_core_server.ServerConfig{
		ProjectID:       "privyfellowship-6a4e1",
		CredentialsFile: "/home/slvr/FirebaseTestingApi/config/privyfellowship-6a4e1-firebase-adminsdk-g0c2d-66e133ed37.json",
	}

	CallGrpc(&config)
}
