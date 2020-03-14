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
)

type Build struct {
}

type Coordinator struct {
	log       *zap.Logger
	mux       *http.ServeMux
	fileCache *filecache.Cache

	mu            sync.Mutex
	scheduledJobs map[build.ID]*scheduledJob
	queue         []*scheduledJob
}

func NewCoordinator(
	log *zap.Logger,
	fileCache *filecache.Cache,
) *Coordinator {
	c := &Coordinator{
		log:       log,
		mux:       http.NewServeMux(),
		fileCache: fileCache,

		scheduledJobs: make(map[build.ID]*scheduledJob),
	}

	c.mux.HandleFunc("/build", c.Build)
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

		s := c.scheduleJob(&job)
		<-s.done

		c.log.Debug("job finished", zap.String("job_id", job.ID.String()))

		update := proto.StatusUpdate{JobFinished: s.finished}
		if err := enc.Encode(update); err != nil {
			return err
		}
	}

	update := proto.StatusUpdate{BuildFinished: &proto.BuildFinished{}}
	return enc.Encode(update)
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

	for _, job := range req.FinishedJob {
		job := job

		scheduled, ok := c.lookupJob(job.ID)
		if !ok {
			continue
		}

		c.log.Debug("job finished")
		scheduled.finish(&job)
	}

	var rsp proto.HeartbeatResponse

	var job *build.Job
	for i := 0; i < 10; i++ {
		var ok bool
		job, ok = c.pickJob()

		if ok {
			rsp.JobsToRun = map[build.ID]proto.JobSpec{
				job.ID: {Job: *job},
			}

			break
		}

		time.Sleep(time.Millisecond)
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
