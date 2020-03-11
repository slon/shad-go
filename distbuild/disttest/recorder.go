package disttest

import (
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type JobResult struct {
	Stdout string
	Stderr string

	Code  *int
	Error string
}

type Recorder struct {
	Jobs map[build.ID]*JobResult
}

func NewRecorder() *Recorder {
	return &Recorder{
		Jobs: map[build.ID]*JobResult{},
	}
}

func (r *Recorder) job(jobID build.ID) *JobResult {
	j, ok := r.Jobs[jobID]
	if !ok {
		j = &JobResult{}
		r.Jobs[jobID] = j
	}
	return j
}

func (r *Recorder) OnJobStdout(jobID build.ID, stdout []byte) error {
	j := r.job(jobID)
	j.Stdout += string(stdout)
	return nil
}

func (r *Recorder) OnJobStderr(jobID build.ID, stderr []byte) error {
	j := r.job(jobID)
	j.Stderr += string(stderr)
	return nil
}

func (r *Recorder) OnJobFinished(jobID build.ID) error {
	j := r.job(jobID)
	j.Code = new(int)
	return nil
}

func (r *Recorder) OnJobFailed(jobID build.ID, code int, error string) error {
	j := r.job(jobID)
	j.Code = &code
	j.Error = error
	return nil
}
