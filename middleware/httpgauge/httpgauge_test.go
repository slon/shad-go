package httpgauge_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/middleware/httpgauge"
)

func TestMiddleware(t *testing.T) {
	g := httpgauge.New()

	m := chi.NewRouter()
	m.Use(g.Wrap)

	m.Get("/simple", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	m.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("bug")
	})
	m.Get("/user/{userID}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/simple", nil))
	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/simple", nil))

	require.Panics(t, func() {
		m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/panic", nil))
	})

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := range 1000 {
				m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", fmt.Sprintf("/user/%d", j), nil))
			}
		}()
	}

	wg.Wait()

	require.Equal(t, map[string]int{
		"/simple":        2,
		"/panic":         1,
		"/user/{userID}": 10000,
	}, g.Snapshot())

	w := httptest.NewRecorder()
	g.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))

	require.Equal(t, "/panic 1\n/simple 2\n/user/{userID} 10000\n", w.Body.String())
}
