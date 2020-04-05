// +build !solution

package api

import (
	"net/http"

	"go.uber.org/zap"
)

type HeartbeatHandler struct {
}

func NewHeartbeatHandler(l *zap.Logger, s HeartbeatService) *HeartbeatHandler {
	panic("implement me")
}

func (h *HeartbeatHandler) Register(mux *http.ServeMux) {
	panic("implement me")
}
