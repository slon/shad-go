package keylock

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkMutex_Baseline(b *testing.B) {
	var mu sync.Mutex

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			_ = 0
			mu.Unlock()
		}
	})
}

func BenchmarkKeyLock_SingleKey(b *testing.B) {
	l := New()

	keys := []string{"a"}

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, unlock := l.LockKeys(keys, nil)
			unlock()
		}
	})
}

func BenchmarkKeyLock_MultipleKeys(b *testing.B) {
	l := New()

	keys := []string{"a", "b", "c", "d"}

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, unlock := l.LockKeys(keys, nil)
			unlock()
		}
	})
}

func BenchmarkKeyLock_DifferentKeys(b *testing.B) {
	l := New()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		keys := []string{strconv.Itoa(rand.Int())}

		for pb.Next() {
			_, unlock := l.LockKeys(keys, nil)
			unlock()
		}
	})
}

func BenchmarkKeyLock_NoBusyWait(b *testing.B) {
	l := New()

	lockedKey := []string{"locked"}
	l.LockKeys(lockedKey, nil)

	cancel := make(chan struct{})
	defer close(cancel)
	for i := 0; i < 1000; i++ {
		go func() {
			l.LockKeys(lockedKey, cancel)
		}()
	}

	b.ResetTimer()

	openKey := []string{"a"}
	for i := 0; i < b.N; i++ {
		canceled, unlock := l.LockKeys(openKey, nil)
		if canceled {
			b.Fatal("spurious lock fail")
		}
		unlock()
	}
}
