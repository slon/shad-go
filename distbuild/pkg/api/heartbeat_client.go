// +build !solution

package api

import (
	"context"

	"go.uber.org/zap"
)

type HeartbeatClient struct {
}

func NewHeartbeatClient(l *zap.Logger, endpoint string) *HeartbeatClient {
	panic("implement me")
}

func (c *HeartbeatClient) Heartbeat(ctx context.Context, req *HeartbeatRequest) (*HeartbeatResponse, error) {
	panic("implement me")
}
