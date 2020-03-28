package scheduler

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"go.uber.org/zap/zaptest"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
)

const (
	workerID0 proto.WorkerID = "w0"
)

func TestScheduler(t *testing.T) {
	defer goleak.VerifyNone(t)

	clock := clockwork.NewFakeClock()
	timeAfter = clock.After
	defer func() { timeAfter = time.After }()

	config := Config{
		CacheTimeout: time.Second,
		DepsTimeout:  time.Minute,
	}

	t.Run("SingleJob", func(t *testing.T) {
		s := NewScheduler(zaptest.NewLogger(t), config)

		job0 := &build.Job{ID: build.NewID()}
		pendingJob0 := s.ScheduleJob(job0)

		s.RegisterWorker(workerID0)
		pickerJob := s.PickJob(workerID0, nil)

		require.Equal(t, pendingJob0, pickerJob)

		result := &proto.JobResult{ID: job0.ID, ExitCode: 0}
		s.OnJobComplete(workerID0, job0.ID, result)

		select {
		case <-pendingJob0.Finished:
			require.Equal(t, pendingJob0.Result, result)

		default:
			t.Fatalf("job0 is not finished")
		}
	})

	t.Run("PickJobTimeout", func(t *testing.T) {
		s := NewScheduler(zaptest.NewLogger(t), config)

		canceled := make(chan struct{})
		close(canceled)

		s.RegisterWorker(workerID0)
		require.Nil(t, s.PickJob(workerID0, canceled))
	})

	t.Run("CacheLocalScheduling", func(t *testing.T) {
		s := NewScheduler(zaptest.NewLogger(t), config)

		job0 := &build.Job{ID: build.NewID()}
		job1 := &build.Job{ID: build.NewID()}

		s.RegisterWorker(workerID0)
		s.OnJobComplete(workerID0, job0.ID, &proto.JobResult{})

		pendingJob1 := s.ScheduleJob(job1)
		pendingJob0 := s.ScheduleJob(job0)

		// job0 scheduling should be blocked on CacheTimeout
		clock.BlockUntil(1)

		pickedJob := s.PickJob(workerID0, nil)
		require.Equal(t, pendingJob0, pickedJob)

		pickedJob = s.PickJob(workerID0, nil)
		require.Equal(t, pendingJob1, pickedJob)

		clock.Advance(time.Hour)
	})

	t.Run("DependencyLocalScheduling", func(t *testing.T) {
		s := NewScheduler(zaptest.NewLogger(t), config)

		job0 := &build.Job{ID: build.NewID()}
		job1 := &build.Job{ID: build.NewID(), Deps: []build.ID{job0.ID}}
		job2 := &build.Job{ID: build.NewID()}

		s.RegisterWorker(workerID0)
		s.OnJobComplete(workerID0, job0.ID, &proto.JobResult{})

		pendingJob2 := s.ScheduleJob(job2)
		pendingJob1 := s.ScheduleJob(job1)

		// job1 should be blocked on DepsTimeout
		clock.BlockUntil(1)

		pickedJob := s.PickJob(workerID0, nil)
		require.Equal(t, pendingJob1, pickedJob)

		pickedJob = s.PickJob(workerID0, nil)
		require.Equal(t, pendingJob2, pickedJob)

		clock.Advance(time.Hour)
	})
}
