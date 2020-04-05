package worker

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	"gitlab.com/slon/shad-go/distbuild/pkg/artifact"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
)

type Worker struct {
	id                  api.WorkerID
	coordinatorEndpoint string

	log *zap.Logger

	fileCache  *filecache.Cache
	fileClient *filecache.Client
	fileOnce   singleflight.Group

	artifacts *artifact.Cache

	mux       *http.ServeMux
	heartbeat *api.HeartbeatClient

	mu           sync.Mutex
	newArtifacts []build.ID
	newSources   []build.ID
	finishedJobs []api.JobResult
}

func New(
	workerID api.WorkerID,
	coordinatorEndpoint string,
	log *zap.Logger,
	fileCache *filecache.Cache,
	artifacts *artifact.Cache,
) *Worker {
	w := &Worker{
		id:                  workerID,
		coordinatorEndpoint: coordinatorEndpoint,
		log:                 log,

		fileCache: fileCache,
		artifacts: artifacts,

		fileClient: filecache.NewClient(log, coordinatorEndpoint),
		heartbeat:  api.NewHeartbeatClient(log, coordinatorEndpoint),

		mux: http.NewServeMux(),
	}

	artifact.NewHandler(w.log, w.artifacts).Register(w.mux)
	return w
}

func (w *Worker) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w.mux.ServeHTTP(rw, r)
}

func (w *Worker) recover() error {
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
		w.log.Debug("sending heartbeat request")
		rsp, err := w.heartbeat.Heartbeat(ctx, w.buildHeartbeat())
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			w.log.DPanic("heartbeat failed", zap.Error(err))
			continue
		}
		w.log.Debug("received heartbeat response",
			zap.Int("num_jobs", len(rsp.JobsToRun)))

		for _, spec := range rsp.JobsToRun {
			spec := spec

			w.log.Debug("running job", zap.String("job_id", spec.Job.ID.String()))
			result, err := w.runJob(ctx, &spec)
			if err != nil {
				errStr := fmt.Sprintf("job %s failed: %v", spec.Job.ID, err)

				w.log.Debug("job failed", zap.String("job_id", spec.Job.ID.String()), zap.Error(err))
				w.jobFinished(&api.JobResult{ID: spec.Job.ID, Error: &errStr})
				continue
			}

			w.log.Debug("job finished", zap.String("job_id", spec.Job.ID.String()))
			w.jobFinished(result)
		}
	}
}
