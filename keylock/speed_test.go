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
