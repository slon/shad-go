package simpleserver

import (
	"net/http"
)

func RunServer() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}

	err := http.ListenAndServe(":8080", http.HandlerFunc(handler))
	if err != nil {
		panic(err)
	}
}

func RunTLSServer() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}

	err := http.ListenAndServeTLS(":8080", "cert.crt", "private.key", http.HandlerFunc(handler))
	if err != nil {
		panic(err)
	}
}
