package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/slon/shad-go/firegod/dontlook"
)

var Solved = false

var flagPort = flag.String("http", "", "")

var router = chi.NewRouter()

func main() {
	flag.Parse()

	dontlook.NewService(router)

	log.Fatal(http.ListenAndServe(*flagPort, router))
}
