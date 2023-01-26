package callgrpc

import (
	"firebaseapi/collectionx"
	"testing"
)

func Test_Call_Grpc(t *testing.T) {
	config := collectionx.Config{
		GrpcAddress:        "0.0.0.0:9090",
		FirebasePath:       "/home/slvr/FirebaseTestingApi/config/privyfellowship-6a4e1-firebase-adminsdk-g0c2d-66e133ed37.json",
		PubSubTopic:        "",
		ExternalCollection: "development-privypass_collection-core-se",
		ProjectID:          "privyfellowship-6a4e1",
	}

	CallGrpc(&config)
}
