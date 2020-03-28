package scheduler

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
)

type PendingJob struct {
	Job      *build.Job
	Result   *proto.JobResult
	Finished chan struct{}

	mu       sync.Mutex
	pickedUp chan struct{}
}

func (p *PendingJob) finish(res *proto.JobResult) {
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

type jobQueue struct {
	mu   sync.Mutex
	jobs []*PendingJob
}

func (q *jobQueue) put(job *PendingJob) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.jobs = append(q.jobs, job)
}

func (q *jobQueue) pop() *PendingJob {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.jobs) == 0 {
		return nil
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job
}

type Config struct {
	CacheTimeout time.Duration
	DepsTimeout  time.Duration
}

type Scheduler struct {
	l      *zap.Logger
	config Config

	mu sync.Mutex

	cachedJobs  map[build.ID]map[proto.WorkerID]struct{}
	pendingJobs map[build.ID]*PendingJob

	cacheLocalQueue map[proto.WorkerID]*jobQueue
	depLocalQueue   map[proto.WorkerID]*jobQueue
	globalQueue     chan *PendingJob
}

func NewScheduler(l *zap.Logger, config Config) *Scheduler {
	return &Scheduler{
		l:      l,
		config: config,

		cachedJobs:  make(map[build.ID]map[proto.WorkerID]struct{}),
		pendingJobs: make(map[build.ID]*PendingJob),

		cacheLocalQueue: make(map[proto.WorkerID]*jobQueue),
		depLocalQueue:   make(map[proto.WorkerID]*jobQueue),
		globalQueue:     make(chan *PendingJob),
	}
}

func (c *Scheduler) RegisterWorker(workerID proto.WorkerID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cacheLocalQueue[workerID]
	if ok {
		return
	}

	c.cacheLocalQueue[workerID] = new(jobQueue)
	c.depLocalQueue[workerID] = new(jobQueue)
}

func (c *Scheduler) OnJobComplete(workerID proto.WorkerID, jobID build.ID, res *proto.JobResult) bool {
	c.l.Debug("job completed", zap.String("worker_id", workerID.String()), zap.String("job_id", jobID.String()))

	c.mu.Lock()
	pendingJob, pendingFound := c.pendingJobs[jobID]
	if pendingFound {
		delete(c.pendingJobs, jobID)
	}

	job, ok := c.cachedJobs[jobID]
	if !ok {
		job = make(map[proto.WorkerID]struct{})
		c.cachedJobs[jobID] = job
	}
	job[workerID] = struct{}{}

	c.mu.Unlock()

	if !pendingFound {
		return false
	}

	c.l.Debug("finishing pending job", zap.String("job_id", jobID.String()))
	pendingJob.finish(res)
	return true
}

func (c *Scheduler) findOptimalWorkers(jobID build.ID, deps []build.ID) (cacheLocal, depLocal []proto.WorkerID) {
	depLocalSet := map[proto.WorkerID]struct{}{}

	c.mu.Lock()
	defer c.mu.Unlock()

	for workerID := range c.cachedJobs[jobID] {
		cacheLocal = append(cacheLocal, workerID)
	}

	for _, dep := range deps {
		for workerID := range c.cachedJobs[dep] {
			if _, ok := depLocalSet[workerID]; !ok {
				depLocal = append(depLocal, workerID)
				depLocalSet[workerID] = struct{}{}
			}
		}
	}

	return
}

var timeAfter = time.After

func (c *Scheduler) doScheduleJob(job *PendingJob) {
	cacheLocal, depLocal := c.findOptimalWorkers(job.Job.ID, job.Job.Deps)

	if len(cacheLocal) != 0 {
		c.mu.Lock()
		for _, workerID := range cacheLocal {
			c.cacheLocalQueue[workerID].put(job)
		}
		c.mu.Unlock()

		c.l.Debug("job is put into cache-local queues", zap.String("job_id", job.Job.ID.String()))
		select {
		case <-job.pickedUp:
			c.l.Debug("job picked", zap.String("job_id", job.Job.ID.String()))
			return
		case <-timeAfter(c.config.CacheTimeout):
		}
	}

	if len(depLocal) != 0 {
		c.mu.Lock()
		for _, workerID := range depLocal {
			c.depLocalQueue[workerID].put(job)
		}
		c.mu.Unlock()

		c.l.Debug("job is put into dep-local queues", zap.String("job_id", job.Job.ID.String()))
		select {
		case <-job.pickedUp:
			c.l.Debug("job picked", zap.String("job_id", job.Job.ID.String()))
			return
		case <-timeAfter(c.config.DepsTimeout):
		}
	}

	c.l.Debug("job is put into global queue", zap.String("job_id", job.Job.ID.String()))
	select {
	case c.globalQueue <- job:
	case <-job.pickedUp:
	}
	c.l.Debug("job picked", zap.String("job_id", job.Job.ID.String()))
}

func (c *Scheduler) ScheduleJob(job *build.Job) *PendingJob {
	c.mu.Lock()
	pendingJob, running := c.pendingJobs[job.ID]
	if !running {
		pendingJob = &PendingJob{
			Job:      job,
			Finished: make(chan struct{}),

			pickedUp: make(chan struct{}),
		}

		c.pendingJobs[job.ID] = pendingJob
	}
	c.mu.Unlock()

	if !running {
		c.l.Debug("job is scheduled", zap.String("job_id", job.ID.String()))
		go c.doScheduleJob(pendingJob)
	} else {
		c.l.Debug("job is pending", zap.String("job_id", job.ID.String()))
	}

	return pendingJob
}

func (c *Scheduler) PickJob(workerID proto.WorkerID, canceled <-chan struct{}) *PendingJob {
	c.l.Debug("picking next job", zap.String("worker_id", workerID.String()))

	var cacheLocal, depLocal *jobQueue

	c.mu.Lock()
	cacheLocal = c.cacheLocalQueue[workerID]
	depLocal = c.depLocalQueue[workerID]
	c.mu.Unlock()

	for {
		job := cacheLocal.pop()
		if job == nil {
			break
		}

		if job.pickUp() {
			c.l.Debug("picked job from cache-local queue", zap.String("worker_id", workerID.String()), zap.String("job_id", job.Job.ID.String()))
			return job
		}
	}

	for {
		job := depLocal.pop()
		if job == nil {
			break
		}

		if job.pickUp() {
			c.l.Debug("picked job from dep-local queue", zap.String("worker_id", workerID.String()), zap.String("job_id", job.Job.ID.String()))
			return job
		}
	}

	for {
		select {
		case job := <-c.globalQueue:
			if job.pickUp() {
				c.l.Debug("picked job from global queue", zap.String("worker_id", workerID.String()), zap.String("job_id", job.Job.ID.String()))
				return job
			}

		case <-canceled:
			return nil
		}
	}
}
