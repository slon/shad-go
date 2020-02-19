package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/tools/testtool"
)

const importPath = "gitlab.com/slon/shad-go/urlshortener"

var binCache testtool.BinCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var teardown testtool.CloseFunc
		binCache, teardown = testtool.NewBinCache()
		defer teardown()

		return m.Run()
	}())
}

func startServer(t *testing.T) (port string, stop func()) {
	binary, err := binCache.GetBinary(importPath)
	require.NoError(t, err)

	port, err = testtool.GetFreePort()
	require.NoError(t, err, "unable to get free port")

	cmd := exec.Command(binary, "-port", port)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	require.NoError(t, cmd.Start())

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	stop = func() {
		_ = cmd.Process.Kill()
		<-done
	}

	if err = testtool.WaitForPort(t, time.Second*5, port); err != nil {
		stop()
	}

	require.NoError(t, err)
	return
}

func add(t *testing.T, c *resty.Client, shortenerURL, request string) string {
	type Response struct {
		URL string `json:"url"`
		Key string `json:"key"`
	}

	resp, err := c.R().
		SetBody(map[string]interface{}{"url": request}).
		SetResult(&Response{}).
		Post(shortenerURL)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode())
	response := resp.Result().(*Response)
	require.Equal(t, request, response.URL)
	require.Contains(t, resp.Header().Get("Content-Type"), "application/json")

	return response.Key
}

func TestURLShortener_redirect(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	var mu sync.Mutex
	redirects := make(map[string]struct{})

	redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		redirects[r.RequestURI] = struct{}{}
		mu.Unlock()
		_, _ = fmt.Fprintln(w, "hello")
	}))

	client := resty.New()
	addURL := fmt.Sprintf("http://localhost:%s/shorten", port)

	requests := make(map[string]struct{})
	for i := 0; i < 10; i++ {
		path := "/" + testtool.RandomName()
		req := redirectTarget.URL + path
		requests[path] = struct{}{}
		key := add(t, client, addURL, req)

		getURL := fmt.Sprintf("http://localhost:%s/go/%s", port, key)
		resp, err := client.R().Get(getURL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())
	}

	mu.Lock()
	defer mu.Unlock()

	require.True(t, reflect.DeepEqual(requests, redirects),
		fmt.Sprintf("expected: %+v, got: %+v", requests, redirects))
}

func TestURLShortener_badRequest(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	u := fmt.Sprintf("http://localhost:%s/shorten", port)
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(`{"url":"abc}`)).
		Post(u)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
}

func TestURLShortener_badKey(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	u := fmt.Sprintf("http://localhost:%s/go/%s", port, testtool.RandomName())
	resp, err := resty.New().
		SetRedirectPolicy(resty.RedirectPolicyFunc(func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		})).R().
		Get(u)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode())
}

func TestURLShortener_consistency(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	client := resty.New()

	get := func(originalURL, key string) {
		getURL := fmt.Sprintf("http://localhost:%s/go/%s", port, key)

		resp, err := client.
			SetRedirectPolicy(resty.RedirectPolicyFunc(func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			})).R().
			Get(getURL)

		require.NoError(t, err)
		require.Equal(t, http.StatusFound, resp.StatusCode())
		require.Contains(t, resp.Header().Get("Location"), originalURL)
	}

	var urls []string
	for i := 0; i < 10; i++ {
		urls = append(urls, testtool.RandomName())
	}

	keyToURL := make(map[string]string)
	urlToKey := make(map[string]string)

	addURL := fmt.Sprintf("http://localhost:%s/shorten", port)
	for _, u := range urls {
		key := add(t, client, addURL, u)

		url, ok := keyToURL[key]
		require.False(t, ok, fmt.Sprintf("duplicate key %s for urls [%s, %s]", key, u, url))

		urlToKey[u] = key
		keyToURL[key] = u
	}

	for _, u := range urls {
		get(u, urlToKey[u])
	}

	for _, u := range urls {
		key := add(t, client, addURL, u)

		url, ok := keyToURL[key]
		require.True(t, ok, fmt.Sprintf("different keys for the same url %s", u))
		require.Equal(t, url, u)
	}
}
