package main

import (
	"net/http"
	_ "net/http/pprof"
)

func main() {
	for i := 0; i < 10000; i++ {
		go func() {
			select {}
		}()
	}

	http.ListenAndServe(":8080", nil)
}
