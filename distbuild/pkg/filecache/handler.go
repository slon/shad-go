// +build !solution

package filecache

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
}

func NewHandler(l *zap.Logger, cache *Cache) *Handler {
	panic("implement me")
}

func (h *Handler) Register(mux *http.ServeMux) {
	panic("implement me")
}
