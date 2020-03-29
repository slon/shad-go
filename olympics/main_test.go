package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/tools/testtool"
)

const importPath = "gitlab.com/slon/shad-go/olympics"

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

	cmd := exec.Command(binary, "-port", port, "-data", "./testdata/olympicWinners.json")
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

func TestServer_valid(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	c := resty.New().SetTimeout(time.Second)

	for _, e := range []string{"athlete-info", "top-athletes-in-sport", "top-countries-in-year"} {
		t.Run(e, func(t *testing.T) {
			testDir := path.Join("./testdata", "tests", e)
			files, err := ioutil.ReadDir(testDir)
			require.NoError(t, err)

			for _, f := range files {
				if !f.IsDir() {
					continue
				}
				if _, err := strconv.Atoi(f.Name()); err != nil {
					continue
				}

				t.Run(f.Name(), func(t *testing.T) {
					in, err := ioutil.ReadFile(path.Join(testDir, f.Name(), "in.json"))
					require.NoError(t, err)

					out, err := ioutil.ReadFile(path.Join(testDir, f.Name(), "out.json"))
					require.NoError(t, err)

					var values map[string]interface{}
					require.NoError(t, json.Unmarshal(in, &values))

					resp, err := c.R().
						SetQueryParams(toURLValues(values)).
						Get(fmt.Sprintf("http://localhost:%s/%s", port, e))

					require.NoError(t, err)
					require.Equal(t, http.StatusOK, resp.StatusCode())
					require.Contains(t, resp.Header().Get("Content-Type"), "application/json")

					var got interface{}
					err = json.Unmarshal(resp.Body(), &got)
					if err != nil {
						t.Fatalf("Could not parse response body: %v", err)
					}

					var want interface{}
					_ = json.Unmarshal(out, &want)

					if diff := cmp.Diff(want, got); diff != "" {
						t.Errorf("Response diff (-want +got):\n%s", diff)
					}
				})
			}
		})
	}
}

func TestServer_invalid(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	c := resty.New().SetTimeout(time.Second)

	for _, tc := range []struct {
		endpoint     string
		description  string
		queryParams  map[string]string
		expectedCode int
	}{
		{
			endpoint:    "athlete-info",
			description: "name-not-found",
			queryParams: map[string]string{
				"name": "AB",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			endpoint:    "top-athletes-in-sport",
			description: "sport-not-found",
			queryParams: map[string]string{
				"sport": "chess",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			endpoint:    "top-athletes-in-sport",
			description: "invalid-limit",
			queryParams: map[string]string{
				"sport": "Canoeing",
				"limit": "2.5",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			endpoint:    "top-countries-in-year",
			description: "year-not-found",
			queryParams: map[string]string{
				"year": "2009",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			endpoint:    "top-countries-in-year",
			description: "invalid-limit",
			queryParams: map[string]string{
				"year":  "2012",
				"limit": "2.5",
			},
			expectedCode: http.StatusBadRequest,
		},
	} {
		t.Run(tc.endpoint+"-"+tc.description, func(t *testing.T) {
			resp, err := c.R().
				SetQueryParams(tc.queryParams).
				Get(fmt.Sprintf("http://localhost:%s/%s", port, tc.endpoint))

			require.NoError(t, err)
			require.Equal(t, tc.expectedCode, resp.StatusCode())
		})
	}
}

func toURLValues(in map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for k, v := range in {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}
