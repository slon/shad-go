package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

func NewBuildService(l *zap.Logger, s Service) *BuildHandler {
	return &BuildHandler{
		l: l,
		s: s,
	}
}

type BuildHandler struct {
	l *zap.Logger
	s Service
}

func (h *BuildHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/build", h.build)
	mux.HandleFunc("/signal", h.signal)
}

type statusWriter struct {
	id      build.ID
	h       *BuildHandler
	written bool
	w       http.ResponseWriter
	flush   http.Flusher
	enc     *json.Encoder
}

func (w *statusWriter) Started(rsp *BuildStarted) error {
	w.id = rsp.ID
	w.written = true

	w.h.l.Debug("build started", zap.String("build_id", w.id.String()), zap.Any("started", rsp))

	w.w.Header().Set("content-type", "application/json")
	w.w.WriteHeader(http.StatusOK)

	defer w.flush.Flush()
	return w.enc.Encode(rsp)
}

func (w *statusWriter) Updated(update *StatusUpdate) error {
	w.h.l.Debug("build updated", zap.String("build_id", w.id.String()), zap.Any("update", update))

	defer w.flush.Flush()
	return w.enc.Encode(update)
}

func (h *BuildHandler) doBuild(w http.ResponseWriter, r *http.Request) error {
	reqJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var req BuildRequest
	if err = json.Unmarshal(reqJSON, &req); err != nil {
		return err
	}

	flush, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("response writer does not implement http.Flusher")
	}

	sw := &statusWriter{h: h, w: w, enc: json.NewEncoder(w), flush: flush}
	err = h.s.StartBuild(r.Context(), &req, sw)

	if err != nil {
		if sw.written {
			_ = sw.Updated(&StatusUpdate{BuildFailed: &BuildFailed{Error: err.Error()}})
			return nil
		}

		return err
	}

	return nil
}

func (h *BuildHandler) build(w http.ResponseWriter, r *http.Request) {
	if err := h.doBuild(w, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "%v", err)
	}
}

func (h *BuildHandler) doSignal(w http.ResponseWriter, r *http.Request) error {
	buildIDParam := r.URL.Query().Get("build_id")
	if buildIDParam == "" {
		return fmt.Errorf(`"build_id" parameter is missing`)
	}

	var buildID build.ID
	if err := buildID.UnmarshalText([]byte(buildIDParam)); err != nil {
		return err
	}

	reqJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var req SignalRequest
	if err = json.Unmarshal(reqJSON, &req); err != nil {
		return err
	}

	rsp, err := h.s.SignalBuild(r.Context(), buildID, &req)
	if err != nil {
		return err
	}

	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		return err
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rspJSON)
	return nil
}

func (h *BuildHandler) signal(w http.ResponseWriter, r *http.Request) {
	if err := h.doSignal(w, r); err != nil {
		h.l.Warn("build signal failed", zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "%v", err)
	}
}
