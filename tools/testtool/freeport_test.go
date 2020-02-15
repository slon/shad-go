package testtool

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetFreePort(t *testing.T) {
	p, err := GetFreePort()
	require.NoError(t, err)
	require.NotEmpty(t, p)
}

func TestWaitForPort(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.Nil(t, err)
	_, port, err := net.SplitHostPort(u.Host)
	require.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	WaitForPort(ctx, port)

	require.NoError(t, ctx.Err())
}

func TestWaitForPort_timeout(t *testing.T) {
	p, err := GetFreePort()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	WaitForPort(ctx, p)

	require.Error(t, ctx.Err())
}
