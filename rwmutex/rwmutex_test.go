package rwmutex

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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
	for i := 0; i < numReaders; i++ {
		go parallelReader(m, clocked, cunlock, cdone)
	}
	// Wait for all parallel RLock()s to succeed.
	for i := 0; i < numReaders; i++ {
		<-clocked
	}
	for i := 0; i < numReaders; i++ {
		cunlock <- true
	}
	// Wait for the goroutines to finish.
	for i := 0; i < numReaders; i++ {
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
	for i := 0; i < numIterations; i++ {
		rwm.RLock()
		n := atomic.AddInt32(activity, 1)
		if n < 1 || n >= 10000 {
			rwm.RUnlock()
			panic(fmt.Sprintf("wlock(%d)\n", n))
		}
		for i := 0; i < 100; i++ {
		}
		atomic.AddInt32(activity, -1)
		rwm.RUnlock()
	}
	cdone <- true
}

func writer(rwm *RWMutex, numIterations int, activity *int32, cdone chan bool) {
	for i := 0; i < numIterations; i++ {
		rwm.Lock()
		n := atomic.AddInt32(activity, 10000)
		if n != 10000 {
			rwm.Unlock()
			panic(fmt.Sprintf("wlock(%d)\n", n))
		}
		for i := 0; i < 100; i++ {
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
	for i = 0; i < numReaders/2; i++ {
		go reader(rwm, numIterations, &activity, cdone)
	}
	go writer(rwm, numIterations, &activity, cdone)
	for ; i < numReaders; i++ {
		go reader(rwm, numIterations, &activity, cdone)
	}
	// Wait for the 2 writers and all readers to finish.
	for i := 0; i < 2+numReaders; i++ {
		<-cdone
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

type CriticalSection struct {
	mu                         sync.Mutex
	readersCount, writersCount int
}

func (cs *CriticalSection) AddToVariable(value *int, count int) {
	cs.mu.Lock()
	*value += count
	cs.mu.Unlock()
}

func (cs *CriticalSection) Reader(t *testing.T, duration time.Duration) {
	cs.AddToVariable(&cs.readersCount, 1)
	cs.Check(t)
	time.Sleep(duration) // do some work
	cs.AddToVariable(&cs.readersCount, -1)
}

func (cs *CriticalSection) Writer(t *testing.T, duration time.Duration) {
	cs.AddToVariable(&cs.writersCount, 1)
	cs.Check(t)
	time.Sleep(duration) // do some work
	cs.AddToVariable(&cs.writersCount, -1)
}

func (cs *CriticalSection) Check(t *testing.T) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if cs.writersCount > 1 {
		t.Errorf("To much writers: %d", cs.writersCount)
	}
	if cs.writersCount == 1 && cs.readersCount > 0 {
		t.Errorf("We have %d readers and %d writers", cs.readersCount, cs.writersCount)
	}
}

func TestAFewReaders(t *testing.T) {
	var wg sync.WaitGroup
	readersCount := 100
	rwm := New()
	cs := new(CriticalSection)
	ch := make(chan struct{})
	wg.Add(readersCount)
	for i := 0; i < readersCount; i++ {
		go func() {
			rwm.RLock()
			cs.Reader(t, 20*time.Millisecond)
			rwm.RUnlock()
			defer wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()
	select {
	case <-ch: //ok
	case <-time.After(25 * time.Millisecond):
		t.Error("too slow, your readers are blocked")
	}
}

func TestAFewWriters(t *testing.T) {
	var wg sync.WaitGroup
	writersCount := 10
	rwm := New()
	cs := new(CriticalSection)
	wg.Add(writersCount)
	for i := 0; i < writersCount; i++ {
		go func() {
			rwm.Lock()
			cs.Writer(t, 10*time.Millisecond)
			rwm.Unlock()
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func TestWriterAfterReaders(t *testing.T) {
	var wg sync.WaitGroup
	rwm := New()
	cs := new(CriticalSection)
	readersCount := 10
	wg.Add(readersCount + 1)
	for i := 0; i < readersCount; i++ {
		go func() {
			rwm.RLock()
			cs.Reader(t, 100*time.Millisecond)
			rwm.RUnlock()
			defer wg.Done()
		}()
	}

	time.Sleep(10 * time.Millisecond)

	go func() {
		rwm.Lock()
		cs.Writer(t, 10*time.Millisecond)
		rwm.Unlock()
		defer wg.Done()
	}()
	wg.Wait()
}

func TestReadersAfterWriters(t *testing.T) {
	var wg sync.WaitGroup
	rwm := New()
	cs := new(CriticalSection)
	RWCount := 10
	wg.Add(2 * RWCount)
	for i := 0; i < RWCount; i++ {
		go func() {
			rwm.Lock()
			cs.Writer(t, 100*time.Millisecond)
			rwm.Unlock()
			defer wg.Done()
		}()
	}

	time.Sleep(20 * time.Millisecond)

	for i := 0; i < RWCount; i++ {
		go func() {
			rwm.RLock()
			cs.Reader(t, 10*time.Millisecond)
			rwm.RUnlock()
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func TestRWStress(t *testing.T) {
	var wg sync.WaitGroup
	rwm := New()
	cs := new(CriticalSection)
	RWCount := 20
	for j := 0; j < 100; j++ {
		wg.Add(2 * RWCount)
		for i := 0; i < RWCount; i++ {
			go func() {
				time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond) // some delay
				rwm.Lock()
				cs.Writer(t, time.Millisecond)
				rwm.Unlock()
				defer wg.Done()
			}()
			go func() {
				time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond) // some delay
				rwm.RLock()
				cs.Reader(t, time.Millisecond)
				rwm.RUnlock()
				defer wg.Done()
			}()
		}
		wg.Wait()
	}
}
