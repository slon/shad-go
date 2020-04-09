package api_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	"gitlab.com/slon/shad-go/distbuild/pkg/api/mock"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

//go:generate mockgen -package mock -destination mock/heartbeat.go . HeartbeatService

func TestHeartbeat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := zaptest.NewLogger(t)
	m := mock.NewMockHeartbeatService(ctrl)
	mux := http.NewServeMux()
	api.NewHeartbeatHandler(l, m).Register(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := api.NewHeartbeatClient(l, server.URL)

	req := &api.HeartbeatRequest{
		WorkerID: "worker0",
	}
	rsp := &api.HeartbeatResponse{
		JobsToRun: map[build.ID]api.JobSpec{
			{0x01}: {Job: build.Job{Name: "cc a.c"}},
		},
	}

	gomock.InOrder(
		m.EXPECT().Heartbeat(gomock.Any(), gomock.Eq(req)).Times(1).Return(rsp, nil),
		m.EXPECT().Heartbeat(gomock.Any(), gomock.Eq(req)).Times(1).Return(nil, fmt.Errorf("build error: foo bar")),
	)

	clientRsp, err := client.Heartbeat(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, rsp, clientRsp)

	_, err = client.Heartbeat(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "build error: foo bar")
}
