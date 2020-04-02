package simpleserver

import (
	"errors"
	"net/http"
)

func RunServerWithMiddleware() {
	getOnly := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}

	err := http.ListenAndServe(":8080", getOnly(http.HandlerFunc(handler)))
	if err != nil {
		panic(err)
	}
}

func UnifiedErrorMiddleware() {
	wrapErrorReply := func(h func(w http.ResponseWriter, r *http.Request) error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := h(w, r); err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
		})
	}

	handler := func(w http.ResponseWriter, r *http.Request) error {
		if r.URL.Query().Get("secret") != "FtP8lu70XjWj8Stt" {
			return errors.New("secret mismatch")
		}
		_, _ = w.Write([]byte("pong"))
		return nil
	}

	err := http.ListenAndServe(":8080", wrapErrorReply(handler))
	if err != nil {
		panic(err)
	}
}
