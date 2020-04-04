package artifact

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/tarstream"
)

type Handler struct {
	l *zap.Logger
	c *Cache
}

func NewHandler(l *zap.Logger, c *Cache) *Handler {
	return &Handler{l: l, c: c}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/artifact", h.artifact)
}

func (h *Handler) doArtifact(w http.ResponseWriter, r *http.Request) error {
	idStr := r.URL.Query().Get("id")

	var id build.ID
	if err := id.UnmarshalText([]byte(idStr)); err != nil {
		return err
	}

	h.l.Debug("streaming artifact", zap.String("artifact_id", id.String()))
	artifactDir, unlock, err := h.c.Get(id)
	if err != nil {
		return err
	}
	defer unlock()

	w.WriteHeader(http.StatusOK)
	if err := tarstream.Send(artifactDir, w); err != nil {
		h.l.Warn("error streaming artifact", zap.Error(err))
	}
	return nil
}

func (h *Handler) artifact(w http.ResponseWriter, r *http.Request) {
	if err := h.doArtifact(w, r); err != nil {
		h.l.Warn("artifact handler error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "%v", err)
	}
}
