package dist

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
	"gitlab.com/slon/shad-go/distbuild/pkg/scheduler"
)

type Coordinator struct {
	log       *zap.Logger
	mux       *http.ServeMux
	fileCache *filecache.Cache

	mu        sync.Mutex
	builds    map[build.ID]*Build
	scheduler *scheduler.Scheduler
}

var defaultConfig = scheduler.Config{
	CacheTimeout: time.Millisecond * 10,
	DepsTimeout:  time.Millisecond * 100,
}

func NewCoordinator(
	log *zap.Logger,
	fileCache *filecache.Cache,
) *Coordinator {
	c := &Coordinator{
		log:       log,
		mux:       http.NewServeMux(),
		fileCache: fileCache,

		builds:    make(map[build.ID]*Build),
		scheduler: scheduler.NewScheduler(log, defaultConfig),
	}

	apiHandler := api.NewBuildService(log, c)
	apiHandler.Register(c.mux)

	heartbeatHandler := api.NewHeartbeatHandler(log, c)
	heartbeatHandler.Register(c.mux)

	fileHandler := filecache.NewHandler(log, c.fileCache)
	fileHandler.Register(c.mux)

	return c
}

func (c *Coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}

func (c *Coordinator) addBuild(b *Build) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.builds[b.ID] = b
}

func (c *Coordinator) removeBuild(b *Build) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.builds, b.ID)
}

func (c *Coordinator) getBuild(id build.ID) *Build {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.builds[id]
}

func (c *Coordinator) StartBuild(ctx context.Context, req *api.BuildRequest, w api.StatusWriter) error {
	b := NewBuild(&req.Graph, c)

	c.addBuild(b)
	defer c.removeBuild(b)

	return b.Run(ctx, w)
}

func (c *Coordinator) SignalBuild(ctx context.Context, buildID build.ID, signal *api.SignalRequest) (*api.SignalResponse, error) {
	b := c.getBuild(buildID)
	if b == nil {
		return nil, fmt.Errorf("build %q not found", buildID)
	}

	return b.Signal(ctx, signal)
}

func (c *Coordinator) Heartbeat(ctx context.Context, req *api.HeartbeatRequest) (*api.HeartbeatResponse, error) {
	c.scheduler.RegisterWorker(req.WorkerID)

	for _, job := range req.FinishedJob {
		job := job

		c.scheduler.OnJobComplete(req.WorkerID, job.ID, &job)
	}

	rsp := &api.HeartbeatResponse{
		JobsToRun: map[build.ID]api.JobSpec{},
	}

	job := c.scheduler.PickJob(req.WorkerID, ctx.Done())
	if job != nil {
		rsp.JobsToRun[job.Job.ID] = api.JobSpec{Job: *job.Job}
	}

	return rsp, nil
}
