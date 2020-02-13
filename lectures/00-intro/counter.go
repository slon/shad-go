package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	counter int
	mu      sync.Mutex
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var i int

	mu.Lock()
	i = counter
	counter++
	mu.Unlock()

	fmt.Fprintf(w, "counter = %d\n", i)
}
