// +build !solution

package api

import (
	"net/http"

	"go.uber.org/zap"
)

func NewBuildService(l *zap.Logger, s Service) *BuildHandler {
	panic("implement me")
}

type BuildHandler struct {
}

func (h *BuildHandler) Register(mux *http.ServeMux) {
	panic("implement me")
}
