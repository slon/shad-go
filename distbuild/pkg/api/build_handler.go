package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

func NewServiceHandler(l *zap.Logger, s Service) *ServiceHandler {
	return &ServiceHandler{
		l: l,
		s: s,
	}
}

type ServiceHandler struct {
	l *zap.Logger
	s Service
}

func (s *ServiceHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/build", s.build)
	mux.HandleFunc("/signal", s.signal)
}

type statusWriter struct {
	written bool
	w       http.ResponseWriter
	enc     *json.Encoder
}

func (w *statusWriter) Started(rsp *BuildStarted) error {
	w.written = true
	w.w.Header().Set("content-type", "application/json")
	w.w.WriteHeader(http.StatusOK)
	return w.enc.Encode(rsp)
}

func (w *statusWriter) Updated(update *StatusUpdate) error {
	return w.enc.Encode(update)
}

func (s *ServiceHandler) doBuild(w http.ResponseWriter, r *http.Request) error {
	reqJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var req BuildRequest
	if err = json.Unmarshal(reqJSON, &req); err != nil {
		return err
	}

	sw := &statusWriter{w: w, enc: json.NewEncoder(w)}
	err = s.s.StartBuild(r.Context(), &req, sw)

	if err != nil {
		if sw.written {
			_ = sw.Updated(&StatusUpdate{BuildFailed: &BuildFailed{Error: err.Error()}})
			return nil
		}

		return err
	}

	return nil
}

func (s *ServiceHandler) build(w http.ResponseWriter, r *http.Request) {
	if err := s.doBuild(w, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "%v", err)
	}
}

func (s *ServiceHandler) doSignal(w http.ResponseWriter, r *http.Request) error {
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

	rsp, err := s.s.SignalBuild(r.Context(), buildID, &req)
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

func (s *ServiceHandler) signal(w http.ResponseWriter, r *http.Request) {
	if err := s.doSignal(w, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "%v", err)
	}
}
