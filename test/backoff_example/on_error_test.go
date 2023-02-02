package backoffexample

import (
	"context"
	"testing"
	"time"

	"github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Some result that the client might return.
type fakeResponse struct{}

// Some client that can perform RPCs.
type fakeClient struct{}

// PerformSomeRPC is a fake RPC that a client might perform.
func (c *fakeClient) PerformSomeRPC(ctx context.Context) (*fakeResponse, error) {
	// An actual client would return something meaningful here.
	return nil, nil
}

func Test_Perform_Some_RPC(t *testing.T) {
	ctx := context.Background()
	c := &fakeClient{}

	shouldRetryUnavailableUnKnown := func(err error) bool {
		st, ok := status.FromError(err)
		if !ok {
			return false
		}

		return st.Code() == codes.Unavailable || st.Code() == codes.Unknown
	}
	retryer := gax.OnErrorFunc(gax.Backoff{
		Initial:    time.Second,
		Max:        32 * time.Second,
		Multiplier: 2,
	}, shouldRetryUnavailableUnKnown)

	performSomeRPCWithRetry := func(ctx context.Context) (*fakeResponse, error) {
		for {
			resp, err := c.PerformSomeRPC(ctx)
			if err != nil {
				if delay, shouldRetry := retryer.Retry(err); shouldRetry {
					if err := gax.Sleep(ctx, delay); err != nil {
						return nil, err
					}
					continue
				}
				return nil, err
			}
			return resp, err
		}
	}

	// It's recommended to set deadlines on RPCs and around retrying. This is
	// also usually preferred over setting some fixed number of retries: one
	// advantage this has is that backoff settings can be changed independently
	// of the deadline, whereas with a fixed number of retries the deadline
	// would be a constantly-shifting goalpost.
	ctxWithTimeout, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Minute))
	defer cancel()

	resp, err := performSomeRPCWithRetry(ctxWithTimeout)
	if err != nil {
		// TODO: handle err
		t.Error(err)
	}
	_ = resp // TODO: use resp if err is nil
}
