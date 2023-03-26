package requestlog_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"gitlab.com/slon/shad-go/middleware/requestlog"
)

func TestRequestLog(t *testing.T) {
	core, obs := observer.New(zap.DebugLevel)

	m := chi.NewRouter()
	m.Use(requestlog.Log(zap.New(core)))

	m.Get("/simple", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	m.Post("/post", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	m.Get("/forbidden", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	m.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 200)
		w.WriteHeader(http.StatusOK)
	})
	m.Get("/forgetful", func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})
	m.Get("/buggy", func(w http.ResponseWriter, r *http.Request) {
		panic("bug")
	})

	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/simple", nil))
	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/post", nil))
	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/forbidden", nil))
	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/slow", nil))
	m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/forgetful", nil))

	require.Panics(t, func() {
		m.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/buggy", nil))
	})

	checkEntries := func(path string, panic bool, code int) {
		entries := obs.FilterField(zap.String("path", path)).All()

		require.Len(t, entries, 2)
		require.Equal(t, entries[0].Message, "request started")
		require.Contains(t, entries[0].ContextMap(), "request_id")

		var requestID zap.Field
		for _, f := range entries[0].Context {
			if f.Key == "request_id" {
				requestID = f
			}
		}

		if !panic {
			require.Equal(t, entries[1].Message, "request finished")
			require.Contains(t, entries[1].Context, zap.Int("status_code", code))
			require.Contains(t, entries[1].ContextMap(), "duration")
			require.Contains(t, entries[1].Context, requestID)
		} else {
			require.Equal(t, entries[1].Message, "request panicked")
			require.Contains(t, entries[1].Context, requestID)
		}
	}

	checkEntries("/simple", false, http.StatusOK)
	checkEntries("/post", false, http.StatusOK)
	checkEntries("/forbidden", false, http.StatusForbidden)
	checkEntries("/slow", false, http.StatusOK)
	checkEntries("/forgetful", false, http.StatusOK)
	checkEntries("/buggy", true, 0)
}
