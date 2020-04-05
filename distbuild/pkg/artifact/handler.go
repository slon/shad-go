// +build !solution

package artifact

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
}

func NewHandler(l *zap.Logger, c *Cache) *Handler {
	panic("implement me")
}

func (h *Handler) Register(mux *http.ServeMux) {
	panic("implement me")
}
