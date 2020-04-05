package scheduler

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type PendingJob struct {
	Job      *api.JobSpec
	Finished chan struct{}
	Result   *api.JobResult

	mu       sync.Mutex
	pickedUp chan struct{}
}

func (p *PendingJob) finish(res *api.JobResult) {
	p.Result = res
	close(p.Finished)
}

func (p *PendingJob) pickUp() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case <-p.pickedUp:
		return false
	default:
		close(p.pickedUp)
		return true
	}
}

func (p *PendingJob) enqueue(q chan *PendingJob) {
	select {
	case q <- p:
	case <-p.pickedUp:
	}
}

type workerQueue struct {
	cacheQueue chan *PendingJob
	depQueue   chan *PendingJob
}

type Config struct {
	CacheTimeout time.Duration
	DepsTimeout  time.Duration
}

type Scheduler struct {
	l      *zap.Logger
	config Config

	mu sync.Mutex

	cachedJobs map[build.ID]map[api.WorkerID]struct{}

	pendingJobs    map[build.ID]*PendingJob
	pendingJobDeps map[build.ID]map[*PendingJob]struct{}

	workerQueue map[api.WorkerID]*workerQueue
	globalQueue chan *PendingJob
}

func NewScheduler(l *zap.Logger, config Config) *Scheduler {
	return &Scheduler{
		l:      l,
		config: config,

		cachedJobs:     make(map[build.ID]map[api.WorkerID]struct{}),
		pendingJobs:    make(map[build.ID]*PendingJob),
		pendingJobDeps: make(map[build.ID]map[*PendingJob]struct{}),

		workerQueue: make(map[api.WorkerID]*workerQueue),
		globalQueue: make(chan *PendingJob),
	}
}

func (c *Scheduler) RegisterWorker(workerID api.WorkerID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.workerQueue[workerID]
	if ok {
		return
	}

	c.workerQueue[workerID] = &workerQueue{
		cacheQueue: make(chan *PendingJob),
		depQueue:   make(chan *PendingJob),
	}
}

func (c *Scheduler) OnJobComplete(workerID api.WorkerID, jobID build.ID, res *api.JobResult) bool {
	c.l.Debug("job completed", zap.String("worker_id", workerID.String()), zap.String("job_id", jobID.String()))

	c.mu.Lock()
	pendingJob, pendingFound := c.pendingJobs[jobID]
	if pendingFound {
		delete(c.pendingJobs, jobID)
	}

	job, ok := c.cachedJobs[jobID]
	if !ok {
		job = make(map[api.WorkerID]struct{})
		c.cachedJobs[jobID] = job
	}
	job[workerID] = struct{}{}

	workerQueue := c.workerQueue[workerID]
	for waiter := range c.pendingJobDeps[jobID] {
		go waiter.enqueue(workerQueue.depQueue)
	}

	c.mu.Unlock()

	if !pendingFound {
		return false
	}

	c.l.Debug("finishing pending job", zap.String("job_id", jobID.String()))
	pendingJob.finish(res)
	return true
}

func (c *Scheduler) enqueueCacheLocal(job *PendingJob) bool {
	cached := false

	for workerID := range c.cachedJobs[job.Job.ID] {
		cached = true
		go job.enqueue(c.workerQueue[workerID].cacheQueue)
	}

	return cached
}

var timeAfter = time.After

func (c *Scheduler) putDepQueue(job *PendingJob, dep build.ID) {
	depJobs, ok := c.pendingJobDeps[dep]
	if !ok {
		depJobs = make(map[*PendingJob]struct{})
		c.pendingJobDeps[dep] = depJobs
	}
	depJobs[job] = struct{}{}
}

func (c *Scheduler) deleteDepQueue(job *PendingJob, dep build.ID) {
	depJobs := c.pendingJobDeps[dep]
	delete(depJobs, job)
	if len(depJobs) == 0 {
		delete(c.pendingJobDeps, dep)
	}
}

func (c *Scheduler) doScheduleJob(job *PendingJob, cached bool) {
	if cached {
		select {
		case <-job.pickedUp:
			c.l.Debug("job picked", zap.String("job_id", job.Job.ID.String()))
			return
		case <-timeAfter(c.config.CacheTimeout):
		}
	}

	c.mu.Lock()
	workers := make(map[api.WorkerID]struct{})

	for _, dep := range job.Job.Deps {
		c.putDepQueue(job, dep)

		for workerID := range c.cachedJobs[dep] {
			if _, ok := workers[workerID]; ok {
				return
			}

			go job.enqueue(c.workerQueue[workerID].depQueue)
			workers[workerID] = struct{}{}
		}
	}
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		for _, dep := range job.Job.Deps {
			c.deleteDepQueue(job, dep)
		}
	}()

	c.l.Debug("job is put into dep-local queues", zap.String("job_id", job.Job.ID.String()))

	select {
	case <-job.pickedUp:
		c.l.Debug("job picked", zap.String("job_id", job.Job.ID.String()))
		return
	case <-timeAfter(c.config.DepsTimeout):
	}

	go job.enqueue(c.globalQueue)
	c.l.Debug("job is put into global queue", zap.String("job_id", job.Job.ID.String()))

	<-job.pickedUp
	c.l.Debug("job picked", zap.String("job_id", job.Job.ID.String()))
}

func (c *Scheduler) ScheduleJob(job *api.JobSpec) *PendingJob {
	var cached bool

	c.mu.Lock()
	pendingJob, running := c.pendingJobs[job.ID]
	if !running {
		pendingJob = &PendingJob{
			Job:      job,
			Finished: make(chan struct{}),

			pickedUp: make(chan struct{}),
		}

		c.pendingJobs[job.ID] = pendingJob
		cached = c.enqueueCacheLocal(pendingJob)
	}
	c.mu.Unlock()

	if !running {
		c.l.Debug("job is scheduled", zap.String("job_id", job.ID.String()))
		go c.doScheduleJob(pendingJob, cached)
	} else {
		c.l.Debug("job is pending", zap.String("job_id", job.ID.String()))
	}

	return pendingJob
}

func (c *Scheduler) PickJob(ctx context.Context, workerID api.WorkerID) *PendingJob {
	c.l.Debug("picking next job", zap.String("worker_id", workerID.String()))

	c.mu.Lock()
	local := c.workerQueue[workerID]
	c.mu.Unlock()

	var pg *PendingJob
	var queue string

	for {
		select {
		case pg = <-c.globalQueue:
			queue = "global"
		case pg = <-local.depQueue:
			queue = "dep"
		case pg = <-local.cacheQueue:
			queue = "cache"
		case <-ctx.Done():
			return nil
		}

		if pg.pickUp() {
			break
		}
	}

	c.l.Debug("picked job",
		zap.String("worker_id", workerID.String()),
		zap.String("job_id", pg.Job.ID.String()),
		zap.String("queue", queue))

	return pg
}
