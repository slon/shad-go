package client

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
)

type Client struct {
	l         *zap.Logger
	client    *api.BuildClient
	cache     *filecache.Client
	sourceDir string
}

func NewClient(
	l *zap.Logger,
	apiEndpoint string,
	sourceDir string,
) *Client {
	return &Client{
		l:         l,
		client:    api.NewBuildClient(l, apiEndpoint),
		cache:     filecache.NewClient(l, apiEndpoint),
		sourceDir: sourceDir,
	}
}

type BuildListener interface {
	OnJobStdout(jobID build.ID, stdout []byte) error
	OnJobStderr(jobID build.ID, stderr []byte) error

	OnJobFinished(jobID build.ID) error
	OnJobFailed(jobID build.ID, code int, error string) error
}

func (c *Client) uploadSources(ctx context.Context, graph *build.Graph, started *api.BuildStarted) error {
	for _, id := range started.MissingFiles {
		c.l.Debug("uploading missing file to coordinator", zap.String("id", id.String()))

		path, ok := graph.SourceFiles[id]
		if !ok {
			return fmt.Errorf("file is missing in build graph: id=%s", id)
		}

		if err := c.cache.Upload(ctx, id, path); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Build(ctx context.Context, graph build.Graph, lsn BuildListener) error {
	started, r, err := c.client.StartBuild(ctx, &api.BuildRequest{Graph: graph})
	if err != nil {
		return err
	}

	c.l.Debug("build started", zap.String("build_id", started.ID.String()))
	if err := c.uploadSources(ctx, &graph, started); err != nil {
		return err
	}

	uploadDone := &api.SignalRequest{UploadDone: &api.UploadDone{}}
	_, err = c.client.SignalBuild(ctx, started.ID, uploadDone)
	if err != nil {
		return err
	}

	for {
		u, err := r.Next()
		if err == io.EOF {
			return fmt.Errorf("unexpected end of status stream")
		} else if err != nil {
			return err
		}

		c.l.Debug("received status update", zap.String("build_id", started.ID.String()), zap.Any("update", u))
		switch {
		case u.BuildFailed != nil:
			return fmt.Errorf("build failed: %s", u.BuildFailed.Error)

		case u.BuildFinished != nil:
			return nil

		case u.JobFinished != nil:
			jf := u.JobFinished

			if jf.Stdout != nil {
				if err := lsn.OnJobStdout(jf.ID, jf.Stdout); err != nil {
					return err
				}
			}

			if jf.Stderr != nil {
				if err := lsn.OnJobStderr(jf.ID, jf.Stderr); err != nil {
					return err
				}
			}

			if jf.Error != nil {
				if err := lsn.OnJobFailed(jf.ID, jf.ExitCode, *jf.Error); err != nil {
					return err
				}
			} else {
				if err := lsn.OnJobFinished(jf.ID); err != nil {
					return err
				}
			}

		default:
			return fmt.Errorf("build failed: unexpected status update")
		}
	}
}
