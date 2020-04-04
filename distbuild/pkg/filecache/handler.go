package filecache

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type Handler struct {
	l      *zap.Logger
	cache  *Cache
	single singleflight.Group
}

func NewHandler(l *zap.Logger, cache *Cache) *Handler {
	return &Handler{
		l:     l,
		cache: cache,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/file", h.file)
}

func (h *Handler) doGet(w http.ResponseWriter, r *http.Request, id build.ID) error {
	path, unlock, err := h.cache.Get(id)
	if err != nil {
		return err
	}
	defer unlock()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(w, f); err != nil {
		h.l.Warn("error streaming file", zap.Error(err))
	}
	return nil
}

func (h *Handler) doPut(w http.ResponseWriter, r *http.Request, id build.ID) error {
	_, err, _ := h.single.Do(id.String(), func() (interface{}, error) {
		w, abort, err := h.cache.Write(id)
		if errors.Is(err, ErrExists) {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		defer abort()

		if _, err = io.Copy(w, r.Body); err != nil {
			return nil, err
		}
		return nil, w.Close()
	})

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *Handler) file(w http.ResponseWriter, r *http.Request) {
	var id build.ID
	err := id.UnmarshalText([]byte(r.URL.Query().Get("id")))

	if err == nil {
		switch r.Method {
		case http.MethodGet:
			err = h.doGet(w, r, id)
		case http.MethodPut:
			err = h.doPut(w, r, id)
		default:
			err = fmt.Errorf("filehandler: unsupported method %s", r.Method)
		}
	}

	if err != nil {
		h.l.Warn("file error", zap.String("method", r.Method), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "%v", err)
	}
}
