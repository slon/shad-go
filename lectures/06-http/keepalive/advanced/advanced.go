package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	urls := []string{"https://golang.org/doc", "https://golang.org/pkg", "https://golang.org/help"}
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := client.Get(url)
			if err != nil {
				fmt.Printf("%s: %s\n", url, err)
				return
			}
			fmt.Printf("%s - %d\n", url, resp.StatusCode)
		}(url)
	}

	wg.Wait()
}
