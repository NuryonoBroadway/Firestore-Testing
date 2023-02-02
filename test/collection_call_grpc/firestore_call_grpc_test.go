package collectioncallgrpc

import (
	collectionxserver "firebaseapi/collectionx/collectionx_server"
	"testing"
)

func Test_Call_Grpc(t *testing.T) {
	config := collectionxserver.ServerConfig{
		FirestoreProjectID:       "privyfellowship-6a4e1",
		FirestoreCredentialsFile: "/home/slvr/FirebaseTestingApi/config/privyfellowship-6a4e1-firebase-adminsdk-g0c2d-66e133ed37.json",
		PubSubProjectID:          "cellular-effect-306806",
		PubSubCredentialsFile:    "/home/slvr/FirebaseTestingApi/config/cellular-effect-306806-afdfaa2f69e4.json",
	}

	CallGrpc(&config)
}
