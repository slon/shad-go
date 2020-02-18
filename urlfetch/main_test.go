// +build !change

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/tools/testtool"
)

const urlfetchImportPath = "gitlab.com/slon/shad-go/urlfetch"

var binCache testtool.BinCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var teardown testtool.CloseFunc
		binCache, teardown = testtool.NewBinCache()
		defer teardown()

		return m.Run()
	}())
}

func TestUrlFetch_valid(t *testing.T) {
	binary, err := binCache.GetBinary(urlfetchImportPath)
	require.NoError(t, err)

	type endpoint string
	type response string

	for _, tc := range []struct {
		name     string
		h        http.HandlerFunc
		queries  []endpoint
		expected []response
	}{
		{
			name: "404",
			h: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "The requested URL was not found.", http.StatusNotFound)
			},
			queries:  []endpoint{"/" + endpoint(testtool.RandomName())},
			expected: []response{"The requested URL was not found."},
		},
		{
			name: "200",
			h: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("The requested URL was found.\n"))
			},
			queries:  []endpoint{"/"},
			expected: []response{"The requested URL was found.\n"},
		},
		{
			name: "multiple-urls",
			h: func(w http.ResponseWriter, r *http.Request) {
				mux := http.NewServeMux()
				mux.HandleFunc("/foo", func(w http.ResponseWriter, h *http.Request) {
					_, _ = w.Write([]byte("foo"))
				})
				mux.HandleFunc("/bar", func(w http.ResponseWriter, h *http.Request) {
					_, _ = w.Write([]byte("bar"))
				})
				mux.ServeHTTP(w, r)
			},
			queries:  []endpoint{"/foo", "/bar"},
			expected: []response{"foo", "bar"},
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

			data, err := cmd.Output()
			require.NoError(t, err)

			for _, r := range tc.expected {
				require.True(t, bytes.Contains(data, []byte(r)),
					fmt.Sprintf(`output="%s" does not contain expected response="%s"`, data, r))
			}

		})
	}
}

func TestUrlFetch_malformed(t *testing.T) {
	binary, err := binCache.GetBinary(urlfetchImportPath)
	require.NoError(t, err)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("success"))
	}))
	defer s.Close()

	for _, tc := range []struct {
		name string
		urls []string
	}{
		{
			name: "invalid-protocol-scheme",
			urls: []string{"golang.org"},
		},
		{
			name: "valid+invalid",
			urls: []string{s.URL, "golang.org"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binary, tc.urls...)
			cmd.Stdout = nil
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				return
			}

			t.Fatalf("process ran with err=%v, want exit status != 0", err)
		})
	}
}

func TestUrlFetch_order(t *testing.T) {
	binary, err := binCache.GetBinary(urlfetchImportPath)
	require.NoError(t, err)

	mux := http.NewServeMux()

	var expectedCallOrder []int

	var mu sync.Mutex
	var callOrder []int

	n := 1000
	for i := 0; i < n; i++ {
		i := i
		expectedCallOrder = append(expectedCallOrder, i)
		s := strconv.Itoa(i)
		mux.HandleFunc("/"+s, func(w http.ResponseWriter, h *http.Request) {
			mu.Lock()
			callOrder = append(callOrder, i)
			mu.Unlock()
			_, _ = w.Write([]byte(s))
		})
	}

	s := httptest.NewServer(mux)
	defer s.Close()

	urls := make([]string, n)
	for i := 0; i < n; i++ {
		urls[i] = s.URL + "/" + strconv.Itoa(i)
	}

	cmd := exec.Command(binary, urls...)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	require.NoError(t, cmd.Run())

	mu.Lock()
	defer mu.Unlock()

	require.Equal(t, expectedCallOrder, callOrder)
}
