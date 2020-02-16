// +build !change

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/tools/testtool"
)

const fetchallImportPath = "gitlab.com/slon/shad-go/fetchall"

var binCache testtool.BinCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var teardown testtool.CloseFunc
		binCache, teardown = testtool.NewBinCache()
		defer teardown()

		return m.Run()
	}())
}

func TestFetchall_valid(t *testing.T) {
	binary, err := binCache.GetBinary(fetchallImportPath)
	require.NoError(t, err)

	type endpoint string

	for _, tc := range []struct {
		name    string
		h       http.HandlerFunc
		queries []endpoint
	}{
		{
			name: "404",
			h: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "The requested URL was not found.", http.StatusNotFound)
			},
			queries: []endpoint{"/" + endpoint(testtool.RandomName())},
		},
		{
			name: "200",
			h: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("The requested URL was found.\n"))
			},
			queries: []endpoint{"/"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := httptest.NewServer(tc.h)
			defer s.Close()

			urls := make([]string, len(tc.queries))
			for i, q := range tc.queries {
				urls[i] = s.URL + string(q)
			}

			cmd := exec.Command(binary, urls...)
			cmd.Stdout = nil
			cmd.Stderr = os.Stderr

			require.NoError(t, cmd.Run())
		})
	}
}

func TestFetchall_multipleURLs(t *testing.T) {
	binary, err := binCache.GetBinary(fetchallImportPath)
	require.NoError(t, err)

	var fooHit, barHit int32

	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, h *http.Request) {
		atomic.StoreInt32(&fooHit, 1)
		_, _ = w.Write([]byte("foo"))
	})
	mux.HandleFunc("/bar", func(w http.ResponseWriter, h *http.Request) {
		atomic.StoreInt32(&barHit, 1)
		_, _ = w.Write([]byte("bar"))
	})

	s := httptest.NewServer(mux)
	defer s.Close()

	cmd := exec.Command(binary, s.URL+"/foo", s.URL+"/bar")
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	require.NoError(t, cmd.Run())

	require.Equal(t, int32(1), atomic.LoadInt32(&fooHit))
	require.Equal(t, int32(1), atomic.LoadInt32(&barHit))
}

func TestFetchall_malformed(t *testing.T) {
	binary, err := binCache.GetBinary(fetchallImportPath)
	require.NoError(t, err)

	hit := int32(0)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hit, 1)
		_, _ = w.Write([]byte("success"))
	}))
	defer s.Close()

	cmd := exec.Command(binary, "golang.org", s.URL, s.URL)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	require.NoError(t, err)

	require.True(t, atomic.LoadInt32(&hit) >= 2)
}

func TestFetchall_concurrency(t *testing.T) {
	binary, err := binCache.GetBinary(fetchallImportPath)
	require.NoError(t, err)

	var mu sync.Mutex
	var callOrder []time.Duration

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := r.URL.Query().Get("duration")
		require.NotEmpty(t, s)

		d, err := time.ParseDuration(s)
		require.NoError(t, err)

		time.Sleep(d)

		mu.Lock()
		callOrder = append(callOrder, d)
		mu.Unlock()

		_, _ = fmt.Fprintln(w, "hello")
	}))
	defer s.Close()

	makeURL := func(d time.Duration) string {
		v := url.Values{}
		v.Add("duration", d.String())
		return fmt.Sprintf("%s?%s", s.URL, v.Encode())
	}

	fastURL := makeURL(time.Millisecond * 10)
	slowURL := makeURL(time.Second)

	cmd := exec.Command(binary, slowURL, fastURL)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	require.NoError(t, cmd.Run())

	mu.Lock()
	defer mu.Unlock()

	require.Len(t, callOrder, 2)
	require.True(t, sort.SliceIsSorted(callOrder, func(i, j int) bool {
		return callOrder[i] < callOrder[j]
	}))
}
