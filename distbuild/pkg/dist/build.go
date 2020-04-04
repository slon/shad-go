package dist

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type Build struct {
	ID    build.ID
	Graph *build.Graph

	reverseFiles map[string]build.ID

	l          *zap.Logger
	c          *Coordinator
	uploadDone chan struct{}
}

func NewBuild(graph *build.Graph, c *Coordinator) *Build {
	id := build.NewID()

	return &Build{
		ID:    id,
		Graph: graph,

		reverseFiles: make(map[string]build.ID),

		l:          c.log.With(zap.String("build_id", id.String())),
		c:          c,
		uploadDone: make(chan struct{}),
	}
}

func (b *Build) missingFiles() []build.ID {
	var files []build.ID

	for id, path := range b.Graph.SourceFiles {
		files = append(files, id)
		b.reverseFiles[path] = id
	}

	return files
}

func (b *Build) Run(ctx context.Context, w api.StatusWriter) error {
	if err := w.Started(&api.BuildStarted{ID: b.ID, MissingFiles: b.missingFiles()}); err != nil {
		return err
	}

	b.l.Debug("waiting for file upload")
	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-b.uploadDone:
	}
	b.l.Debug("file upload completed")

	for _, job := range b.Graph.Jobs {
		job := job

		s := b.c.scheduler.ScheduleJob(&job)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.Finished:
		}

		b.l.Debug("job finished", zap.String("job_id", job.ID.String()))

		jobFinished := api.StatusUpdate{JobFinished: s.Result}
		if err := w.Updated(&jobFinished); err != nil {
			return err
		}
	}

	finished := api.StatusUpdate{BuildFinished: &api.BuildFinished{}}
	return w.Updated(&finished)
}

func (b *Build) Signal(ctx context.Context, req *api.SignalRequest) (*api.SignalResponse, error) {
	switch {
	case req.UploadDone != nil:
		select {
		case <-b.uploadDone:
			return nil, fmt.Errorf("upload already done")
		default:
			close(b.uploadDone)
		}

	default:
		return nil, fmt.Errorf("unexpected signal kind")
	}

	return &api.SignalResponse{}, nil
}
