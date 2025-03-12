package keylock

import (
	"math/rand"
	"strconv"
	"testing"
)

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
	for range 1000 {
		go func() {
			l.LockKeys(lockedKey, cancel)
		}()
	}

	

	openKey := []string{"a"}
	for b.Loop() {
		canceled, unlock := l.LockKeys(openKey, nil)
		if canceled {
			b.Fatal("spurious lock fail")
		}
		unlock()
	}
}
