package once

import (
	"testing"

	"gitlab.com/slon/shad-go/tools/testtool"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, once *Once, o *one, c chan bool) {
	once.Do(func() { o.Increment() })
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestOnce(t *testing.T) {
	o := new(one)
	once := New()
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(t, once, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("once failed outside run: %d is not 1", *o)
	}
}

func TestOncePanic(t *testing.T) {
	once := New()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Once.Do did not panic")
			}
		}()
		once.Do(func() {
			panic("failed")
		})
	}()

	once.Do(func() {
		t.Fatalf("Once.Do called twice")
	})
}

func TestOnceManyTimes(t *testing.T) {
	const N = 1000
	for i := 0; i < N; i++ {
		TestOnce(t)
	}
}

func TestOnceNoBusyWait(t *testing.T) {
	once := New()

	done := make(chan struct{})
	defer close(done)

	for i := 0; i < 100; i++ {
		go once.Do(func() {
			<-done
		})
	}

	testtool.VerifyNoBusyGoroutines(t)
}

func TestNoSyncPackageImported(t *testing.T) {
	testtool.CheckForbiddenImport(t, "sync")
}
