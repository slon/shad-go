package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/tools/testtool"
)

const importPath = "gitlab.com/slon/shad-go/firewall/cmd/firewall"

var binCache testtool.BinCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var teardown testtool.CloseFunc
		binCache, teardown = testtool.NewBinCache()
		defer teardown()

		return m.Run()
	}())
}

func storeConfig(t *testing.T, conf string) (filename string, cleanup func()) {
	t.Helper()

	filename = path.Join(os.TempDir(), testtool.RandomName()+".yaml")
	err := ioutil.WriteFile(filename, []byte(conf), 0777)
	require.NoError(t, err)

	cleanup = func() { _ = os.Remove(filename) }
	return
}

func startServer(t *testing.T, serviceURL string, conf string) (port string, stop func()) {
	binary, err := binCache.GetBinary(importPath)
	require.NoError(t, err)

	confPath, removeConf := storeConfig(t, conf)
	defer removeConf()

	port, err = testtool.GetFreePort()
	require.NoError(t, err, "unable to get free port")

	addr := fmt.Sprintf("localhost:%s", port)

	cmd := exec.Command(binary, "-service-addr", serviceURL, "-addr", addr, "-conf", confPath)
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

func TestFirewall(t *testing.T) {
	echoService := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, r.Body)
		defer func() { _ = r.Body.Close() }()
		w.Header().Set("Content-Length", fmt.Sprintf("%d", r.ContentLength))
	}

	c := resty.New()

	type result struct {
		code int
		body string
	}

	for _, tc := range []struct {
		name        string
		conf        string
		service     http.HandlerFunc
		makeRequest func() *resty.Request
		endpoint    string
		expected    result
	}{
		{
			name:    "empty",
			conf:    ``,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody("hello")
			},
			expected: result{code: http.StatusOK, body: "hello"},
		},
		{
			name: "simple",
			conf: `
rules:
  - endpoint: "/"
    forbidden_user_agents:
      - 'python-requests.*'
    forbidden_headers:
      - 'Content-Type: text/html'
    required_headers:
      - "Content-Type"
    max_request_length_bytes: 50
    max_response_length_bytes: 50
    forbidden_response_codes: [201]
    forbidden_request_re:
      - '.*(\.\./){3,}.*'
    forbidden_response_re:
      - '.*admin.*'
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().
					SetBody(`{"user_id": 123, "path": "../../user"}`).
					SetHeaders(map[string]string{
						"User-Agent":   "Mozilla/5.0",
						"Content-Type": "application/json",
					})
			},
			expected: result{code: http.StatusOK, body: `{"user_id": 123, "path": "../../user"}`},
		},
		{
			name: "unprotected-endpoint",
			conf: `
rules:
  - endpoint: "/list"
    forbidden_user_agents:
      - 'python-requests.*'
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().
					SetHeader("User-Agent", "python-requests/2.22.0").
					SetBody(`{"user_id": 123}`)
			},
			endpoint: "/login",
			expected: result{code: http.StatusOK, body: `{"user_id": 123}`},
		},
		{
			name: "empty-rule",
			conf: `
rules:
  - endpoint: "/list"
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().
					SetBody(`{"user_id": 123}`)
			},
			endpoint: "/list",
			expected: result{code: http.StatusOK, body: `{"user_id": 123}`},
		},
		{
			name: "bad-user-agent",
			conf: `
rules:
  - endpoint: "/"
    forbidden_user_agents:
      - 'python-requests.*'
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody("hello").SetHeader("User-Agent", "python-requests/2.22.0")
			},
			expected: result{code: http.StatusForbidden, body: "Forbidden"},
		},
		{
			name: "forbidden-header",
			conf: `
rules:
  - endpoint: "/"
    forbidden_headers:
      - 'Content-Type: text/html'
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody("hello").SetHeader("Content-Type", "text/html")
			},
			expected: result{code: http.StatusForbidden, body: "Forbidden"},
		},
		{
			name: "missing-required-header",
			conf: `
rules:
  - endpoint: "/"
    required_headers:
      - "Content-Type"
      - "Content-Length"
      - "NoOneUsesThisHeader"
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody("hello").SetHeader("Content-Length", "5")
			},
			expected: result{code: http.StatusForbidden, body: "Forbidden"},
		},
		{
			name: "max-request-length-exceeded",
			conf: `
rules:
  - endpoint: "/"
    max_request_length_bytes: 4
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody("hello")
			},
			expected: result{code: http.StatusForbidden, body: "Forbidden"},
		},
		{
			name: "max-response-length-exceeded",
			conf: `
rules:
  - endpoint: "/"
    max_response_length_bytes: 4
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody("hello")
			},
			expected: result{code: http.StatusForbidden, body: "Forbidden"},
		},
		{
			name: "bad-status-code",
			conf: `
rules:
  - endpoint: "/"
    forbidden_response_codes: [200]
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody("hello")
			},
			expected: result{code: http.StatusForbidden, body: "Forbidden"},
		},
		{
			name: "bad-request",
			conf: `
rules:
  - endpoint: "/"
    forbidden_request_re:
      - '.*(\.\./){3,}.*'
`,
			service: func(w http.ResponseWriter, r *http.Request) {},
			makeRequest: func() *resty.Request {
				return c.R().SetBody(`{"path": "../../../../etc.passwd"}`)
			},
			expected: result{code: http.StatusForbidden, body: "Forbidden"},
		},
		{
			name: "bad-response",
			conf: `
rules:
  - endpoint: "/"
    forbidden_response_re:
      - '.*admin.*'
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody(`{"user": "admin", "password": "1234"}`)
			},
			expected: result{code: http.StatusForbidden, body: "Forbidden"},
		},
		{
			name: "many-rules-forbidden",
			conf: `
rules:
  - endpoint: "/list"
    forbidden_response_re:
      - '.*admin.*'
  - endpoint: "/dump"
    max_response_length_bytes: 4
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody(`hello`)
			},
			endpoint: "/dump",
			expected: result{code: http.StatusForbidden, body: "Forbidden"},
		},
		{
			name: "many-rules-ok",
			conf: `
rules:
  - endpoint: "/list"
    forbidden_response_re:
      - '.*admin.*'
  - endpoint: "/dump"
    max_response_length_bytes: 4
`,
			service: echoService,
			makeRequest: func() *resty.Request {
				return c.R().SetBody(`hello`)
			},
			endpoint: "/list",
			expected: result{code: http.StatusOK, body: "hello"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := httptest.NewServer(tc.service)
			defer service.Close()

			port, cleanup := startServer(t, service.URL, tc.conf)
			defer cleanup()

			u := fmt.Sprintf("http://localhost:%s%s", port, tc.endpoint)

			resp, err := tc.makeRequest().Post(u)
			require.NoError(t, err)

			require.Equal(t, tc.expected.code, resp.StatusCode())
			require.Equal(t, tc.expected.body, resp.String())
		})
	}
}
