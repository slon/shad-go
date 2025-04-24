package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"
)

var g any

// START OMIT
func Alloc() {
	for {
		g = make([]byte, 1024*1024)
		time.Sleep(time.Second)
	}
}

func main() {
	for i := 0; i < 10000; i++ {
		go func() {
			select {}
		}()
	}
	go Alloc()
	http.ListenAndServe(":8080", nil)
}

// END OMIT
