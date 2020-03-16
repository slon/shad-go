package keylock_test

import (
	"sync"
	"testing"
	"time"

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

	canceled, unlock := l.LockKeys([]string{"a", "b"}, nil)
	require.False(t, canceled)

	canceled, _ = l.LockKeys([]string{"", "b", "c"}, timeout(time.Millisecond*10))
	require.True(t, canceled)

	unlock()

	canceled, _ = l.LockKeys([]string{"", "b", "c"}, nil)
	require.False(t, canceled)
}

func TestKeyLock_Progress(t *testing.T) {
	defer goleak.VerifyNone(t)
	l := keylock.New()

	canceled, unlock := l.LockKeys([]string{"a", "b"}, nil)
	require.False(t, canceled)
	defer unlock()

	go func() {
		_, _ = l.LockKeys([]string{"b", "c"}, nil)
	}()

	time.Sleep(time.Millisecond * 10)
	canceled, _ = l.LockKeys([]string{"d"}, nil)
	require.False(t, canceled)
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

			for i := 0; i < N; i++ {
				cancelled, unlock := l.LockKeys([]string{"a"}, timeout(time.Millisecond))
				if !cancelled {
					unlock()
				}
			}
		}()
	}

	wg.Wait()
}
