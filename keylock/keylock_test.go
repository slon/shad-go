package keylock_test

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"gitlab.com/slon/shad-go/keylock"
)

func timeout(d time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(d)
		close(ch)
	}()
	return ch
}

func TestKeyLock_Simple(t *testing.T) {
	defer goleak.VerifyNone(t)

	l := keylock.New()

	canceled, unlock0 := l.LockKeys([]string{"a", "b"}, nil)
	require.False(t, canceled)

	canceled, _ = l.LockKeys([]string{"", "b", "c"}, timeout(time.Millisecond*10))
	require.True(t, canceled)

	unlock0()

	canceled, unlock1 := l.LockKeys([]string{"", "b", "c"}, nil)
	require.False(t, canceled)
	unlock1()
}

func TestKeyLock_Progress(t *testing.T) {
	defer goleak.VerifyNone(t)
	l := keylock.New()

	canceled, unlock0 := l.LockKeys([]string{"a", "b"}, nil)
	require.False(t, canceled)
	defer unlock0()

	go func() {
		_, unlock := l.LockKeys([]string{"b", "c"}, nil)
		unlock()
	}()

	time.Sleep(time.Millisecond * 10)
	canceled, unlock1 := l.LockKeys([]string{"d"}, nil)
	require.False(t, canceled)
	unlock1()
}

func TestKeyLock_DeadlockFree(t *testing.T) {
	const N = 10000

	defer goleak.VerifyNone(t)
	l := keylock.New()

	var wg sync.WaitGroup
	wg.Add(3)

	checkLock := func(keys []string) {
		defer wg.Done()

		for i := 0; i < N; i++ {
			cancelled, unlock := l.LockKeys(keys, nil)
			if cancelled {
				t.Error("spurious lock failure")
				return
			}
			unlock()
		}
	}

	go checkLock([]string{"a", "b", "c"})
	go checkLock([]string{"b", "c", "a"})
	go checkLock([]string{"c", "a", "b"})

	wg.Wait()
}

func TestKeyLock_NoMutates(t *testing.T) {
	defer goleak.VerifyNone(t)
	l := keylock.New()

	keys := []string{"b", "c", "a"}
	passedKeys := make([]string, len(keys))
	copy(passedKeys, keys)

	_, unlock := l.LockKeys(passedKeys, nil)
	unlock()

	require.Equal(t, keys, passedKeys, "passed keys shouldn't be mutated")
}

func TestKeyLock_SingleKeyStress(t *testing.T) {
	const (
		N = 1000
		G = 100
	)

	defer goleak.VerifyNone(t)
	l := keylock.New()

	var wg sync.WaitGroup
	wg.Add(G)

	for i := 0; i < G; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < N; j++ {
				cancelled, unlock := l.LockKeys([]string{"a"}, timeout(time.Millisecond))
				if !cancelled {
					unlock()
				}
			}
		}()
	}

	wg.Wait()
}

func TestKeyLock_MutualExclusionStress(t *testing.T) {
	const (
		N = 1000
		G = 100
		M = 15
		K = 3
	)

	defer goleak.VerifyNone(t)

	locked := map[string]bool{}
	var mu sync.Mutex

	l := keylock.New()

	var wg sync.WaitGroup
	wg.Add(G)

	for i := 0; i < G; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < N; j++ {
				keys := []string{}
				for k := 0; k < K; k++ {
					keys = append(keys, fmt.Sprint(rand.Intn(N)))
				}

				slices.Sort(keys)
				keys = slices.Compact(keys)

				_, unlock := l.LockKeys(keys, nil)
				mu.Lock()
				for _, key := range keys {
					assert.False(t, locked[key])
					locked[key] = true
				}
				mu.Unlock()

				time.Sleep(time.Millisecond)

				mu.Lock()
				for _, key := range keys {
					locked[key] = false
				}
				mu.Unlock()

				unlock()
			}
		}()
	}

	wg.Wait()
}
