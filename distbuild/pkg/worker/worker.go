package worker

import (
	"context"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/artifact"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
)

type Worker struct {
	coordinatorEndpoint string

	log *zap.Logger

	fileCache *filecache.Cache
	artifacts *artifact.Cache

	mux *http.ServeMux

	mu           sync.Mutex
	newArtifacts []build.ID
	newSources   []build.ID
}

func New(
	coordinatorEndpoint string,
	log *zap.Logger,
	fileCache *filecache.Cache,
	artifacts *artifact.Cache,
) *Worker {
	return &Worker{
		coordinatorEndpoint: coordinatorEndpoint,
		log:                 log,
		fileCache:           fileCache,
		artifacts:           artifacts,

		mux: http.NewServeMux(),
	}
}

func (w *Worker) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w.mux.ServeHTTP(rw, r)
}

func (w *Worker) recover() error {
	err := w.fileCache.Range(func(file build.ID) error {
		w.newSources = append(w.newSources, file)
		return nil
	})
	if err != nil {
		return err
	}

	return w.artifacts.Range(func(file build.ID) error {
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
