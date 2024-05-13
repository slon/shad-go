package rsem

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/redisfixture"
	"gitlab.com/slon/shad-go/tools/testtool"
	"go.uber.org/goleak"
)

func TestSemaphore_Simple(t *testing.T) {
	goleak.VerifyNone(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisfixture.StartRedis(t),
	})
	defer func() { _ = rdb.Close() }()
	sem := NewSemaphore(rdb)
	ctx := context.Background()

	release, err := sem.Acquire(ctx, "simple", 1)
	require.NoError(t, err)
	require.NoError(t, release())

	release, err = sem.Acquire(ctx, "simple", 1)
	require.NoError(t, err)
	require.NoError(t, release())
}

func TestSemaphore_Limit1(t *testing.T) {
	goleak.VerifyNone(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisfixture.StartRedis(t),
	})
	defer func() { _ = rdb.Close() }()
	sem := NewSemaphore(rdb)
	ctx := context.Background()

	release, err := sem.Acquire(ctx, "limit1", 1)
	require.NoError(t, err)
	defer func() { require.NoError(t, release()) }()

	acquired := make(chan struct{})
	defer func() { <-acquired }()

	go func() {
		defer close(acquired)

		releaseInner, err := sem.Acquire(ctx, "limit1", 1)
		assert.NoError(t, err)
		require.NoError(t, releaseInner())
	}()

	select {
	case <-acquired:
		t.Errorf("semaphore not working")
	case <-time.After(time.Second * 5):
		require.NoError(t, release())
	}
}

func TestSemaphore_IndependentKeys(t *testing.T) {
	goleak.VerifyNone(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisfixture.StartRedis(t),
	})
	defer func() { _ = rdb.Close() }()

	ctx := context.Background()

	for i := 0; i < 1000; i++ {
		sem := NewSemaphore(rdb)

		release, err := sem.Acquire(ctx, fmt.Sprint(i), 1)
		require.NoError(t, err)
		defer func() { require.NoError(t, release()) }()
	}
}

func TestSemaphore_LimitN(t *testing.T) {
	goleak.VerifyNone(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisfixture.StartRedis(t),
	})
	defer func() { _ = rdb.Close() }()
	sem := NewSemaphore(rdb)
	ctx := context.Background()

	const N = 3
	const G = 10

	const testDuration = time.Second * 5
	const lockDuration = time.Millisecond * 100

	startTime := time.Now()

	var counter atomic.Int32
	var wg sync.WaitGroup

	wg.Add(G)
	for g := 0; g < G; g++ {
		go func() {
			defer wg.Done()

			for time.Since(startTime) < testDuration {
				release, err := sem.Acquire(ctx, "limitN", N)
				if !assert.NoError(t, err) {
					return
				}

				counter.Add(1)

				if k := counter.Load(); k > N {
					counter.Add(-1)
					require.NoError(t, release())

					t.Errorf("%d goroutines in critical section", k)
					return
				}

				time.Sleep(lockDuration)

				counter.Add(-1)
				require.NoError(t, release())
			}
		}()
	}

	wg.Wait()
}

func TestSemaphore_ContextCancel(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisfixture.StartRedis(t),
	})
	defer func() { _ = rdb.Close() }()
	sem := NewSemaphore(rdb)
	ctx := context.Background()

	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, err := sem.Acquire(cancelCtx, "cancel", 1)
	require.NoError(t, err)

	cancel()

	release, err := sem.Acquire(ctx, "cancel", 1)
	require.NoError(t, err)

	require.NoError(t, release())
}

var binCache testtool.BinCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var teardown testtool.CloseFunc
		binCache, teardown = testtool.NewBinCache()
		defer teardown()

		return m.Run()
	}())
}

func TestSemaphore_DeadCleanup(t *testing.T) {
	addr := redisfixture.StartRedis(t)

	rdb := redis.NewClient(&redis.Options{Addr: addr})
	defer func() { _ = rdb.Close() }()
	sem := NewSemaphore(rdb)
	ctx := context.Background()

	spawnWorker := func() *exec.Cmd {
		binary, err := binCache.GetBinary("gitlab.com/slon/shad-go/rsem/worker")
		require.NoError(t, err)

		p := exec.Command(binary, addr)
		p.Stderr = os.Stderr

		require.NoError(t, p.Start())

		return p
	}

	checkLocked := func() {
		ctx1, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()

		_, err := sem.Acquire(ctx1, "dead", 2)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	}

	p1 := spawnWorker()
	p2 := spawnWorker()
	time.Sleep(time.Second / 2)

	checkLocked()

	require.NoError(t, p2.Process.Kill())

	release, err := sem.Acquire(ctx, "dead", 2)
	require.NoError(t, err)
	defer func() { _ = release() }()

	checkLocked()

	require.NoError(t, p1.Process.Kill())

	release1, err := sem.Acquire(ctx, "dead", 2)
	require.NoError(t, err)
	defer func() { _ = release1() }()

	checkLocked()
}
