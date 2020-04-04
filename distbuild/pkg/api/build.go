package api

import (
	"context"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type BuildRequest struct {
	Graph build.Graph
}

type BuildStarted struct {
	ID           build.ID
	MissingFiles []build.ID
}

type StatusUpdate struct {
	JobFinished   *JobResult
	BuildFailed   *BuildFailed
	BuildFinished *BuildFinished
}

type BuildFailed struct {
	Error string
}

type BuildFinished struct {
}

type UploadDone struct{}

type SignalRequest struct {
	UploadDone *UploadDone
}

type SignalResponse struct {
}

type StatusWriter interface {
	Started(rsp *BuildStarted) error
	Updated(update *StatusUpdate) error
}

type Service interface {
	StartBuild(ctx context.Context, request *BuildRequest, w StatusWriter) error
	SignalBuild(ctx context.Context, buildID build.ID, signal *SignalRequest) (*SignalResponse, error)
}

type StatusReader interface {
	Close() error
	Next() (*StatusUpdate, error)
}
