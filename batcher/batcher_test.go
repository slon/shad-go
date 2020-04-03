package batcher

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"gitlab.com/slon/shad-go/batcher/slow"
)

func TestSimple(t *testing.T) {
	defer goleak.VerifyNone(t)

	var value slow.Value
	b := NewBatcher(&value)

	value.Store(1)
	require.Equal(t, 1, b.Load())
	require.Equal(t, 1, value.Load())

	value.Store(2)
	require.Equal(t, 2, b.Load())
	require.Equal(t, 2, value.Load())
}

func TestTwoParallelLoads(t *testing.T) {
	defer goleak.VerifyNone(t)
	var value slow.Value
	b := NewBatcher(&value)

	value.Store(1)
	go func() {
		require.Equal(t, 1, b.Load())
	}()
	require.Equal(t, 1, b.Load())
}

func TestStaleRead(t *testing.T) {
	defer goleak.VerifyNone(t)

	const (
		N = 100
		K = 100
		M = 10
	)

	var value slow.Value
	b := NewBatcher(&value)

	var counter int32
	value.Store(counter)

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			time.Sleep(time.Duration(i) * time.Millisecond / time.Duration(N))
			for j := 0; j < K; j++ {
				counterValue := atomic.LoadInt32(&counter)
				batcherValue := b.Load().(int32)

				if batcherValue < counterValue {
					t.Errorf("load returned old value: counter=%d, batcher=%d", counterValue, batcherValue)
					return
				}
			}
		}(i)
	}

	for i := 0; i < M*K; i++ {
		// value is always greater than counter
		value.Store(int32(i))
		atomic.StoreInt32(&counter, int32(i))

		time.Sleep(time.Millisecond / M)
	}

	wg.Wait()
}

func TestSpeed(t *testing.T) {
	defer goleak.VerifyNone(t)

	const (
		N = 100
		K = 200
	)

	var value slow.Value
	b := NewBatcher(&value)

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := 0; i < K; i++ {
				b.Load()
			}
		}()
	}
	wg.Wait()

	require.Truef(t, time.Since(start) < time.Second, "batching it too slow")
}
