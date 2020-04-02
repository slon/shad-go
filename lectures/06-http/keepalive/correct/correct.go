package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

func main() {
	urls := []string{"https://golang.org/doc", "https://golang.org/pkg", "https://golang.org/help"}

	client := &http.Client{Transport: &http.Transport{MaxConnsPerHost: 100}}

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
			defer resp.Body.Close()
			_, _ = io.Copy(ioutil.Discard, resp.Body)
			fmt.Printf("%s - %d\n", url, resp.StatusCode)
		}(url)
	}

	wg.Wait()
}
