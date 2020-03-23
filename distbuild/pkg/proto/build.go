package proto

import (
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type MissingSources struct {
	MissingFiles []build.ID
}

type StatusUpdate struct {
	JobFinished   *JobResult
	BuildFailed   *BuildFailed
	BuildFinished *BuildFinished
}

type BuildFailed struct {
	Error string
}

type BuildFinished struct {
}
