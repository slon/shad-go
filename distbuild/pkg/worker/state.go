package worker

import "gitlab.com/slon/shad-go/distbuild/pkg/api"

func (w *Worker) buildHeartbeat() *api.HeartbeatRequest {
	w.mu.Lock()
	defer w.mu.Unlock()

	req := &api.HeartbeatRequest{
		WorkerID:    w.id,
		FinishedJob: w.finishedJobs,
	}

	w.finishedJobs = nil
	return req
}

func (w *Worker) jobFinished(job *api.JobResult) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.finishedJobs = append(w.finishedJobs, *job)
}
