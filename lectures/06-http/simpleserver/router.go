package simpleserver

import (
	"net/http"
)

func RunServerWithRouting() {
	router := func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/pong":
			pongHandler(w, r)
		case "/shmong":
			shmongHandler(w, r)
		default:
			w.WriteHeader(404)
		}
	}
	err := http.ListenAndServe(":8080", http.HandlerFunc(router))
	if err != nil {
		panic(err)
	}
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}
func shmongHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("shmong"))
}

// OMIT
