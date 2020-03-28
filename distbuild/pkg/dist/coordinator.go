package dist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
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

	c.mux.HandleFunc("/build", c.Build)
	c.mux.HandleFunc("/signal", c.Signal)
	c.mux.HandleFunc("/heartbeat", c.Heartbeat)
	return c
}

func (c *Coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}

func (c *Coordinator) doBuild(w http.ResponseWriter, r *http.Request) error {
	graphJS, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var g build.Graph
	if err := json.Unmarshal(graphJS, &g); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(proto.MissingSources{}); err != nil {
		return err
	}

	for _, job := range g.Jobs {
		job := job

		s := c.scheduler.ScheduleJob(&job)

		select {
		case <-r.Context().Done():
			return r.Context().Err()
		case <-s.Finished:
		}

		c.log.Debug("job finished", zap.String("job_id", job.ID.String()))

		update := proto.StatusUpdate{JobFinished: s.Result}
		if err := enc.Encode(update); err != nil {
			return err
		}
	}

	update := proto.StatusUpdate{BuildFinished: &proto.BuildFinished{}}
	return enc.Encode(update)
}

func (c *Coordinator) Signal(w http.ResponseWriter, r *http.Request) {
	c.log.Debug("build signal started")
	if err := c.doHeartbeat(w, r); err != nil {
		c.log.Error("build signal failed", zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	c.log.Debug("build signal finished")
}

func (c *Coordinator) Build(w http.ResponseWriter, r *http.Request) {
	if err := c.doBuild(w, r); err != nil {
		c.log.Error("build failed", zap.Error(err))

		errorUpdate := proto.StatusUpdate{BuildFailed: &proto.BuildFailed{Error: err.Error()}}
		errorJS, _ := json.Marshal(errorUpdate)
		_, _ = w.Write(errorJS)
	}
}

func (c *Coordinator) doHeartbeat(w http.ResponseWriter, r *http.Request) error {
	var req proto.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	c.scheduler.RegisterWorker(req.WorkerID)

	for _, job := range req.FinishedJob {
		job := job

		c.scheduler.OnJobComplete(req.WorkerID, job.ID, &job)
	}

	rsp := proto.HeartbeatResponse{
		JobsToRun: map[build.ID]proto.JobSpec{},
	}

	job := c.scheduler.PickJob(req.WorkerID, r.Context().Done())
	if job != nil {
		rsp.JobsToRun[job.Job.ID] = proto.JobSpec{Job: *job.Job}
	}

	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		return err
	}

	return nil
}

func (c *Coordinator) Heartbeat(w http.ResponseWriter, r *http.Request) {
	c.log.Debug("heartbeat started")
	if err := c.doHeartbeat(w, r); err != nil {
		c.log.Error("heartbeat failed", zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	c.log.Debug("heartbeat finished")
}
