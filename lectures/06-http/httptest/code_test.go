package httptest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReposCount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("42"))
	}))

	defer srv.Close()

	client := NewAPICLient(srv.URL)
	count, err := client.GetReposCount(context.Background(), "007")

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expectedCount := 42
	if count != expectedCount {
		t.Errorf("expected count to be: %d, got: %d", expectedCount, count)
	}
}

// OMIT
