package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("https://golang.org")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)
}
