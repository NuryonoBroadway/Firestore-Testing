package collectiontest

import (
	collectionxclient "firebaseapi/collectionx/collectionx_client"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Collection_Client_Config(t *testing.T) {
	testCases := []struct {
		name      string
		execution func(t *testing.T)
	}{
		{
			name: "success",
			execution: func(t *testing.T) {
				cfg, err := collectionxclient.NewClientConfig(
					collectionxclient.WithGrpcAddress("0.0.0.0:9090"),
					collectionxclient.WithProjectRootCollection("root_collection"),
					collectionxclient.WithProjectRootDocuments("root_documents"),
					collectionxclient.WithPubSubTopic("pusub"),
				)

				require.NoError(t, err)
				require.NotNil(t, cfg)
			},
		}, {
			name: "error",
			execution: func(t *testing.T) {
				cfg, err := collectionxclient.NewClientConfig(
					collectionxclient.WithGrpcAddress("0.0.0.0:9090"),
					collectionxclient.WithProjectRootDocuments("root_documents"),
					collectionxclient.WithPubSubTopic("pusub"),
				)

				require.Error(t, err)
				require.Nil(t, cfg)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.execution(t)
		})
	}
}
