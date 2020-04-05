package worker

import (
	"context"
	"errors"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	"gitlab.com/slon/shad-go/distbuild/pkg/artifact"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
)

func (w *Worker) downloadFiles(ctx context.Context, files map[build.ID]string) error {
	for id := range files {
		_, unlock, err := w.fileCache.Get(id)
		if errors.Is(err, filecache.ErrNotFound) {
			if err = w.fileClient.Download(ctx, w.fileCache, id); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			unlock()
		}
	}

	return nil
}

func (w *Worker) downloadArtifacts(ctx context.Context, artifacts map[build.ID]api.WorkerID) error {
	for id, worker := range artifacts {
		_, unlock, err := w.artifacts.Get(id)
		if errors.Is(err, artifact.ErrNotFound) {
			if err = artifact.Download(ctx, worker.String(), w.artifacts, id); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			unlock()
		}
	}

	return nil
}
