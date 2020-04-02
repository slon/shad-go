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
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
}
