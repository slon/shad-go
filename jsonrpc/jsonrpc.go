// +build !solution

package jsonrpc

import (
	"context"
	"net/http"
)

func MakeHandler(service interface{}) http.Handler {
	panic("implement me")
}

func Call(ctx context.Context, endpoint string, method string, req, rsp interface{}) error {
	panic("implement me")
}
