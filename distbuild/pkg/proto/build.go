package proto

import (
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type MissingSources struct {
	MissingFiles []build.ID
}

type StatusUpdate struct {
	JobFinished *FinishedJob
	BuildFailed *BuildFailed
}

type BuildFailed struct {
	Error string
}
