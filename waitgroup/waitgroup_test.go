package waitgroup

import (
	"sync/atomic"
	"testing"

	"gitlab.com/slon/shad-go/tools/testtool"
)

func testWaitGroup(t *testing.T, wg1 *WaitGroup, wg2 *WaitGroup) {
	n := 16
	wg1.Add(n)
	wg2.Add(n)
	exited := make(chan bool, n)
	for i := 0; i != n; i++ {
		go func() {
			wg1.Done()
			wg2.Wait()
			exited <- true
		}()
	}
	wg1.Wait()
	for i := 0; i != n; i++ {
		select {
		case <-exited:
			t.Fatal("WaitGroup released group too soon")
		default:
		}
		wg2.Done()
	}
	for i := 0; i != n; i++ {
		<-exited // Will block if barrier fails to unlock someone.
	}
}

func TestWaitGroup(t *testing.T) {
	wg1 := New()
	wg2 := New()

	// Run the same test a few times to ensure barrier is in a proper state.
	for i := 0; i != 8; i++ {
		testWaitGroup(t, wg1, wg2)
	}
}

func recoverFromNegativeCounterPanic(t *testing.T) {
	err := recover()
	if err != "negative WaitGroup counter" {
		t.Fatalf("Unexpected panic: %#v", err)
	}
}

func TestNoop(t *testing.T) {
	wg1 := New()
	wg1.Wait()

	wg1.Add(1)
	go func() {
		wg1.Done()
	}()
	wg1.Wait()

	wg1.Wait()
}

func TestWaitGroupDoneMisuse(t *testing.T) {
	defer recoverFromNegativeCounterPanic(t)
	wg := New()
	wg.Add(1)
	wg.Done()
	wg.Done()
	t.Fatal("Should panic")
}

func TestWaitGroupAddMisuse(t *testing.T) {
	defer recoverFromNegativeCounterPanic(t)
	wg := New()
	wg.Add(1)
	wg.Add(-2)
	t.Fatal("Should panic")
}

func TestWaitGroupRace(t *testing.T) {
	// Run this test for about 1ms.
	for i := 0; i < 1000; i++ {
		wg := New()
		n := new(int32)
		// spawn goroutine 1
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Done()
		}()
		// spawn goroutine 2
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Done()
		}()
		// Wait for goroutine 1 and 2
		wg.Wait()
		if atomic.LoadInt32(n) != 2 {
			t.Fatal("Spurious wakeup from Wait")
		}
	}
}

func TestWaitGroupNoBusyWait(t *testing.T) {
	wg := New()
	wg.Add(1)
	defer wg.Done()

	for i := 0; i < 10; i++ {
		go func() {
			wg.Wait()
		}()
	}

	testtool.VerifyNoBusyGoroutines(t)
}

func TestNoSyncPackageImported(t *testing.T) {
	testtool.CheckForbiddenImport(t, "sync")
}
