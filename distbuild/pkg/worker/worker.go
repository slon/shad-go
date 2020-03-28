package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/artifact"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
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
	finishedJobs []proto.JobResult
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
	//err := w.fileCache.Range(func(file build.ID) error {
	//	w.newSources = append(w.newSources, file)
	//	return nil
	//})
	//if err != nil {
	//	return err
	//}

	return w.artifacts.Range(func(file build.ID) error {
		w.newArtifacts = append(w.newArtifacts, file)
		return nil
	})
}

func (w *Worker) sendHeartbeat(req *proto.HeartbeatRequest) (*proto.HeartbeatResponse, error) {
	reqJS, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", w.coordinatorEndpoint+"/heartbeat", bytes.NewBuffer(reqJS))
	if err != nil {
		return nil, err
	}

	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if httpRsp.StatusCode != http.StatusOK {
		errorString, _ := ioutil.ReadAll(httpRsp.Body)
		return nil, fmt.Errorf("heartbeat failed: %s", errorString)
	}

	var rsp proto.HeartbeatResponse
	if err := json.NewDecoder(httpRsp.Body).Decode(&rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}

func (w *Worker) Run(ctx context.Context) error {
	if err := w.recover(); err != nil {
		return err
	}

	for {
		w.log.Debug("sending heartbeat request")
		rsp, err := w.sendHeartbeat(w.buildHeartbeat())
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
			result, err := w.runJob(ctx, &spec)
			if err != nil {
				errStr := fmt.Sprintf("job %s failed: %v", spec.Job.ID, err)
				w.jobFinished(&proto.JobResult{ID: spec.Job.ID, Error: &errStr})
				continue
			}

			w.jobFinished(result)
		}
	}
}
