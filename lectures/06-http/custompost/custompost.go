package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	body := bytes.NewBufferString("All your base are belong to us")
	req, err := http.NewRequest("POST", "https://myapi.com/create", body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-Source", "Zero Wing")

	repr, err := httputil.DumpRequestOut(req, true)
	if err == nil {
		fmt.Println(string(repr))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
}
