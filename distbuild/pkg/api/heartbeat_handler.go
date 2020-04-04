package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type HeartbeatHandler struct {
	l *zap.Logger
	s HeartbeatService
}

func NewHeartbeatHandler(l *zap.Logger, s HeartbeatService) *HeartbeatHandler {
	return &HeartbeatHandler{l: l, s: s}
}

func (h *HeartbeatHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/heartbeat", h.heartbeat)
}

func (h *HeartbeatHandler) doHeartbeat(w http.ResponseWriter, r *http.Request) error {
	reqJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var req HeartbeatRequest
	if err := json.Unmarshal(reqJSON, &req); err != nil {
		return err
	}

	h.l.Debug("heartbeat started", zap.Any("req", req))
	rsp, err := h.s.Heartbeat(r.Context(), &req)
	if err != nil {
		return err
	}
	h.l.Debug("heartbeat finished", zap.Any("rsp", rsp))

	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rspJSON)
	return nil
}

func (h *HeartbeatHandler) heartbeat(w http.ResponseWriter, r *http.Request) {
	if err := h.doHeartbeat(w, r); err != nil {
		h.l.Warn("heartbeat error", zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "%v", err)
	}
}
