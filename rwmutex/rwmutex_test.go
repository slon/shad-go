package rwmutex

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"gitlab.com/slon/shad-go/tools/testtool"
)

func parallelReader(m *RWMutex, clocked, cunlock, cdone chan bool) {
	m.RLock()
	clocked <- true
	<-cunlock
	m.RUnlock()
	cdone <- true
}

func doTestParallelReaders(numReaders, gomaxprocs int) {
	runtime.GOMAXPROCS(gomaxprocs)
	m := New()
	clocked := make(chan bool)
	cunlock := make(chan bool)
	cdone := make(chan bool)
	for range numReaders {
		go parallelReader(m, clocked, cunlock, cdone)
	}
	// Wait for all parallel RLock()s to succeed.
	for range numReaders {
		<-clocked
	}
	for range numReaders {
		cunlock <- true
	}
	// Wait for the goroutines to finish.
	for range numReaders {
		<-cdone
	}
}

func TestParallelReaders(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(-1))
	doTestParallelReaders(1, 4)
	doTestParallelReaders(3, 4)
	doTestParallelReaders(4, 2)
}

func reader(rwm *RWMutex, numIterations int, activity *int32, cdone chan bool) {
	for range numIterations {
		rwm.RLock()
		n := atomic.AddInt32(activity, 1)
		if n < 1 || n >= 10000 {
			rwm.RUnlock()
			panic(fmt.Sprintf("wlock(%d)\n", n))
		}
		for range 100 {
		}
		atomic.AddInt32(activity, -1)
		rwm.RUnlock()
	}
	cdone <- true
}

func writer(rwm *RWMutex, numIterations int, activity *int32, cdone chan bool) {
	for range numIterations {
		rwm.Lock()
		n := atomic.AddInt32(activity, 10000)
		if n != 10000 {
			rwm.Unlock()
			panic(fmt.Sprintf("wlock(%d)\n", n))
		}
		for range 100 {
		}
		atomic.AddInt32(activity, -10000)
		rwm.Unlock()
	}
	cdone <- true
}

func HammerRWMutex(gomaxprocs, numReaders, numIterations int) {
	runtime.GOMAXPROCS(gomaxprocs)
	// Number of active readers + 10000 * number of active writers.
	var activity int32
	rwm := New()
	cdone := make(chan bool)
	go writer(rwm, numIterations, &activity, cdone)
	var i int
	for i = range numReaders/2 {
		go reader(rwm, numIterations, &activity, cdone)
	}
	go writer(rwm, numIterations, &activity, cdone)
	for ; i < numReaders; i++ {
		go reader(rwm, numIterations, &activity, cdone)
	}
	// Wait for the 2 writers and all readers to finish.
	for range 2+numReaders {
		<-cdone
	}
}

func TestRWMutexReadWrite(t *testing.T) {
	done := make(chan bool)
	go func() {
		rwm := New()
		rwm.RLock()
		rwm.Lock()
		done <- true
	}()

	select {
	case <-time.After(time.Second):
	case <-done:
		t.Fatal("Test finished, must be deadlock")
	}
}

func TestRWMutex(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(-1))
	n := 1000
	if testing.Short() {
		n = 5
	}
	HammerRWMutex(1, 1, n)
	HammerRWMutex(1, 3, n)
	HammerRWMutex(1, 10, n)
	HammerRWMutex(4, 1, n)
	HammerRWMutex(4, 3, n)
	HammerRWMutex(4, 10, n)
	HammerRWMutex(10, 1, n)
	HammerRWMutex(10, 3, n)
	HammerRWMutex(10, 10, n)
	HammerRWMutex(10, 5, n)
}

func TestWriteWriteReadDeadlock(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(-1))
	runtime.GOMAXPROCS(2)
	// Number of active readers + 10000 * number of active writers.
	var activity int32
	rwm := New()
	cdone := make(chan bool, 3)

	for range 2e6 {
		go writer(rwm, 1, &activity, cdone)
		go writer(rwm, 1, &activity, cdone)
		go reader(rwm, 1, &activity, cdone)
		<-cdone
		<-cdone
		<-cdone
	}

}

func TestNoBusyWaitInRlock(t *testing.T) {
	rwm := New()
	rwm.Lock()
	defer rwm.Unlock()

	for i := 0; i < 100; i++ {
		go func() {
			rwm.RLock()
			defer rwm.RUnlock()
		}()
	}

	testtool.VerifyNoBusyGoroutines(t)
}

func TestNoBusyWaitInlock(t *testing.T) {
	rwm := New()
	rwm.RLock()
	defer rwm.RUnlock()

	for i := 0; i < 100; i++ {
		go func() {
			rwm.Lock()
			defer rwm.Unlock()
		}()
	}

	testtool.VerifyNoBusyGoroutines(t)
}

func TestNoSyncPackageImported(t *testing.T) {
	testtool.CheckForbiddenImport(t, "sync")
}
