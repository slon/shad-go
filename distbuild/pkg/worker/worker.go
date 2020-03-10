package worker

import (
	"context"
	"sync"

	"gitlab.com/slon/shad-go/distbuild/pkg/artifact"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
)

type Worker struct {
	CoordinatorEndpoint string

	SourceFiles *filecache.Cache
	Artifacts   *artifact.Cache

	mu           sync.Mutex
	newArtifacts []build.ID
	newSources   []build.ID
}

func (w *Worker) recover() error {
	err := w.SourceFiles.Range(func(file build.ID) error {
		w.newSources = append(w.newSources, file)
		return nil
	})
	if err != nil {
		return err
	}

	return w.Artifacts.Range(func(file build.ID) error {
		w.newArtifacts = append(w.newArtifacts, file)
		return nil
	})
}

func (w *Worker) Run(ctx context.Context) error {
	if err := w.recover(); err != nil {
		return err
	}

	for {

	}
}
