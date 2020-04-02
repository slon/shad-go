// +build !change

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "", "port to listen")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, r.Body)
		defer func() { _ = r.Body.Close() }()
		w.Header().Set("Content-Length", fmt.Sprintf("%d", r.ContentLength))
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
