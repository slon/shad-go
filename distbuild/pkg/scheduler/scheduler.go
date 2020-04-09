// +build !solution

package scheduler

import (
	"context"
	"time"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

var timeAfter = time.After

type PendingJob struct {
	Job      *api.JobSpec
	Finished chan struct{}
	Result   *api.JobResult
}

type Config struct {
	CacheTimeout time.Duration
	DepsTimeout  time.Duration
}

type Scheduler struct {
}

func NewScheduler(l *zap.Logger, config Config) *Scheduler {
	panic("implement me")
}

func (c *Scheduler) LocateArtifact(id build.ID) (api.WorkerID, bool) {
	panic("implement me")
}

func (c *Scheduler) RegisterWorker(workerID api.WorkerID) {
	panic("implement me")
}

func (c *Scheduler) OnJobComplete(workerID api.WorkerID, jobID build.ID, res *api.JobResult) bool {
	panic("implement me")
}

func (c *Scheduler) ScheduleJob(job *api.JobSpec) *PendingJob {
	panic("implement me")
}

func (c *Scheduler) PickJob(ctx context.Context, workerID api.WorkerID) *PendingJob {
	panic("implement me")
}
