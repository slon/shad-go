package httpserver

import (
	"net/http/httptest"
	"testing"
)

func TestHandlerServeHTTP(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	h := handler{}
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected HTTP 200, got: %d", w.Code)
	}
}
