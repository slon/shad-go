package worker

import (
	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
)

func (w *Worker) buildHeartbeat() *proto.HeartbeatRequest {
	w.mu.Lock()
	defer w.mu.Unlock()

	req := &proto.HeartbeatRequest{
		FinishedJob: w.finishedJobs,
	}

	w.finishedJobs = nil
	return req
}

func (w *Worker) jobFinished(job *proto.FinishedJob) {
	w.log.Debug("job finished", zap.String("job_id", job.ID.String()))

	w.mu.Lock()
	defer w.mu.Unlock()

	w.finishedJobs = append(w.finishedJobs, *job)
}
