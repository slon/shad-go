package jsonrpc

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type testService struct{}

type PingRequest struct{}
type PingResponse struct{}

func (*testService) Ping(ctx context.Context, req *PingRequest) (*PingResponse, error) {
	return &PingResponse{}, nil
}

type AddRequest struct{ A, B int }
type AddResponse struct{ Sum int }

func (*testService) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	return &AddResponse{Sum: req.A + req.B}, nil
}

type ErrorRequest struct{}
type ErrorResponse struct{}

func (*testService) Error(ctx context.Context, req *ErrorRequest) (*ErrorResponse, error) {
	return nil, fmt.Errorf("cache is empty")
}

func TestJSONRPC(t *testing.T) {
	server := httptest.NewServer(MakeHandler(&testService{}))
	defer server.Close()

	ctx := context.Background()

	t.Run("Ping", func(t *testing.T) {
		var (
			req PingRequest
			rsp PingResponse
		)

		require.NoError(t, Call(ctx, server.URL, "Ping", &req, &rsp))
	})

	t.Run("Add", func(t *testing.T) {
		var (
			req = AddRequest{A: 1, B: 2}
			rsp AddResponse
		)

		require.NoError(t, Call(ctx, server.URL, "Add", &req, &rsp))
		require.Equal(t, 3, rsp.Sum)
	})

	t.Run("Error", func(t *testing.T) {
		var (
			req ErrorRequest
			rsp ErrorResponse
		)

		err := Call(ctx, server.URL, "Error", &req, &rsp)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cache is empty")
	})
}
