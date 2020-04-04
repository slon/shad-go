package worker

import (
	"context"
	"errors"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
)

func (w *Worker) pullFiles(ctx context.Context, files map[build.ID]string) error {
	for id := range files {
		_, unlock, err := w.fileCache.Get(id)
		if errors.Is(err, filecache.ErrNotFound) {
			if err := w.fileClient.Download(ctx, w.fileCache, id); err != nil {
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
