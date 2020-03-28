package worker

import (
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
)

func (w *Worker) buildHeartbeat() *proto.HeartbeatRequest {
	w.mu.Lock()
	defer w.mu.Unlock()

	req := &proto.HeartbeatRequest{
		WorkerID:    w.id,
		FinishedJob: w.finishedJobs,
	}

	w.finishedJobs = nil
	return req
}

func (w *Worker) jobFinished(job *proto.JobResult) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.finishedJobs = append(w.finishedJobs, *job)
}
