//go:build !solution

package jsonrpc

import (
	"context"
	"net/http"
)

func MakeHandler(service any) http.Handler {
	panic("implement me")
}

func Call(ctx context.Context, endpoint string, method string, req, rsp any) error {
	panic("implement me")
}
