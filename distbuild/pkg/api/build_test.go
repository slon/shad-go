package api_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	mock "gitlab.com/slon/shad-go/distbuild/pkg/api/mock"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

//go:generate mockgen -package mock -destination mock/mock.go . Service

type env struct {
	ctrl   *gomock.Controller
	mock   *mock.MockService
	server *httptest.Server
	client *api.BuildClient
}

func (e *env) stop() {
	e.server.Close()
	e.ctrl.Finish()
}

func newEnv(t *testing.T) (*env, func()) {
	env := &env{}
	env.ctrl = gomock.NewController(t)
	env.mock = mock.NewMockService(env.ctrl)

	log := zaptest.NewLogger(t)

	mux := http.NewServeMux()

	handler := api.NewBuildService(log, env.mock)
	handler.Register(mux)

	env.server = httptest.NewServer(mux)

	env.client = api.NewBuildClient(log, env.server.URL)

	return env, env.stop
}

func TestBuildSignal(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()

	ctx := context.Background()

	buildIDa := build.ID{01}
	buildIDb := build.ID{02}
	req := &api.SignalRequest{}
	rsp := &api.SignalResponse{}

	env.mock.EXPECT().SignalBuild(gomock.Any(), buildIDa, req).Return(rsp, nil)
	env.mock.EXPECT().SignalBuild(gomock.Any(), buildIDb, req).Return(nil, fmt.Errorf("foo bar error"))

	_, err := env.client.SignalBuild(ctx, buildIDa, req)
	require.NoError(t, err)

	_, err = env.client.SignalBuild(ctx, buildIDb, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "foo bar error")
}

func TestBuildStartError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()

	ctx := context.Background()

	env.mock.EXPECT().StartBuild(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo bar error"))

	_, _, err := env.client.StartBuild(ctx, &api.BuildRequest{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "foo bar error")
}

func TestBuildRunning(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()

	ctx := context.Background()

	buildID := build.ID{02}

	req := &api.BuildRequest{
		Graph: build.Graph{SourceFiles: map[build.ID]string{{01}: "a.txt"}},
	}

	started := &api.BuildStarted{ID: buildID}
	finished := &api.StatusUpdate{BuildFinished: &api.BuildFinished{}}

	env.mock.EXPECT().StartBuild(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, req *api.BuildRequest, w api.StatusWriter) error {
			if err := w.Started(started); err != nil {
				return err
			}

			if err := w.Updated(finished); err != nil {
				return err
			}

			return fmt.Errorf("foo bar error")
		})

	rsp, r, err := env.client.StartBuild(ctx, req)
	require.NoError(t, err)
	defer r.Close()

	require.Equal(t, started, rsp)

	u, err := r.Next()
	require.NoError(t, err)
	require.Equal(t, finished, u)

	u, err = r.Next()
	require.NoError(t, err)
	require.Contains(t, u.BuildFailed.Error, "foo bar error")

	_, err = r.Next()
	require.Equal(t, err, io.EOF)
}

func TestBuildResultsStreaming(t *testing.T) {
	// Test is hanging?
	// See https://golang.org/pkg/net/http/#Flusher

	env, stop := newEnv(t)
	defer stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buildID := build.ID{02}
	req := &api.BuildRequest{}
	started := &api.BuildStarted{ID: buildID}

	env.mock.EXPECT().StartBuild(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, req *api.BuildRequest, w api.StatusWriter) error {
			if err := w.Started(started); err != nil {
				return err
			}

			<-ctx.Done()
			return ctx.Err()
		})

	rsp, r, err := env.client.StartBuild(ctx, req)
	require.NoError(t, err)
	defer r.Close()
	require.Equal(t, started, rsp)
}
