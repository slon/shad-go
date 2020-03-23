package dist

import (
	"sync"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
)

type scheduledJob struct {
	job      *build.Job
	finished *proto.JobResult

	mu   sync.Mutex
	done chan struct{}
}

func newScheduledJob(job *build.Job) *scheduledJob {
	return &scheduledJob{
		job:  job,
		done: make(chan struct{}),
	}
}

func (s *scheduledJob) finish(f *proto.JobResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.finished == nil {
		s.finished = f
		close(s.done)
	}
}

func (c *Coordinator) scheduleJob(job *build.Job) *scheduledJob {
	c.mu.Lock()
	defer c.mu.Unlock()

	if scheduled, ok := c.scheduledJobs[job.ID]; ok {
		return scheduled
	} else {
		scheduled = newScheduledJob(job)
		c.scheduledJobs[job.ID] = scheduled
		c.queue = append(c.queue, scheduled)
		return scheduled
	}
}

func (c *Coordinator) pickJob() (*build.Job, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.queue) == 0 {
		return nil, false
	}

	job := c.queue[0].job
	c.queue = c.queue[1:]
	return job, true
}

func (c *Coordinator) lookupJob(id build.ID) (*scheduledJob, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	scheduled, ok := c.scheduledJobs[id]
	return scheduled, ok
}
