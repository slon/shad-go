//go:build !change

package ciletters

type Notification struct {
	Project  GitlabProject
	Branch   string
	Commit   Commit
	Pipeline Pipeline
}

type GitlabProject struct {
	GroupID string
	ID      string
}

type Commit struct {
	// Hash is a 20-byte SHA-1 encoded in hex.
	Hash    string
	Message string
	Author  string
}

type PipelineStatus string

const (
	PipelineStatusOK     PipelineStatus = "ok"
	PipelineStatusFailed PipelineStatus = "failed"
)

type Pipeline struct {
	Status      PipelineStatus
	ID          int64
	TriggeredBy string
	FailedJobs  []Job
}

type Job struct {
	ID        int64
	Name      string
	Stage     string
	RunnerLog string
}
