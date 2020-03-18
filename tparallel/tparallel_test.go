package tparallel

import (
	"sync"
	"testing"
)

type ConcurrencyChecker struct {
	t *testing.T

	mu        sync.Mutex
	nextStage int
	waitCount int
	barrier   chan struct{}
}

func (c *ConcurrencyChecker) Sequential(stage int) {
	c.t.Helper()
	c.t.Logf("Sequential(%d)", stage)

	c.mu.Lock()
	defer c.mu.Unlock()

	if stage != c.nextStage {
		c.t.Errorf("testing method is executed out of sequence: expected=%d, got=%d", c.nextStage, stage)
		return
	}

	c.nextStage++
}

func (c *ConcurrencyChecker) Parallel(stage, count int) {
	c.t.Helper()
	c.t.Logf("Parallel(%d, %d)", stage, count)

	var barrier chan struct{}

	func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		if stage != c.nextStage {
			c.t.Errorf("testing method is executed out of sequence: expected=%d, got=%d", c.nextStage, stage)
			return
		}

		if c.waitCount == 0 {
			c.barrier = make(chan struct{})
		}
		barrier = c.barrier

		c.waitCount++

		if c.waitCount == count {
			c.waitCount = 0
			c.nextStage++
			close(c.barrier)
		}
	}()

	<-barrier
}

func (c *ConcurrencyChecker) Finish(total int) {
	if total != c.nextStage {
		c.t.Errorf("wrong number of stages executed: expected=%d, got=%d", total, c.nextStage)
	}
}

func TestSequentialExecution(t *testing.T) {
	check := &ConcurrencyChecker{t: t}
	defer check.Finish(3)

	Run([]func(*T){
		func(t *T) {
			check.Sequential(0)
		},
		func(t *T) {
			check.Sequential(1)
		},
		func(t *T) {
			check.Sequential(2)
		},
	})
}

func TestParallelExecution(t *testing.T) {
	check := &ConcurrencyChecker{t: t}
	defer check.Finish(4)

	Run([]func(*T){
		func(t *T) {
			check.Sequential(0)
			t.Parallel()
			check.Parallel(3, 3)
		},
		func(t *T) {
			check.Sequential(1)
			t.Parallel()
			check.Parallel(3, 3)
		},
		func(t *T) {
			check.Sequential(2)
			t.Parallel()
			check.Parallel(3, 3)
		},
	})
}

func TestSequentialSubTests(t *testing.T) {
	check := &ConcurrencyChecker{t: t}
	defer check.Finish(5)

	Run([]func(*T){
		func(t *T) {
			check.Sequential(0)

			t.Run(func(t *T) {
				check.Sequential(1)

				t.Run(func(t *T) {
					check.Sequential(2)
				})
			})

			t.Run(func(t *T) {
				check.Sequential(3)
			})
		},
		func(t *T) {
			check.Sequential(4)
		},
	})
}

func TestParallelGroup(t *testing.T) {
	check := &ConcurrencyChecker{t: t}
	defer check.Finish(17)

	Run([]func(*T){
		func(t *T) {
			check.Sequential(0)
		},
		func(t *T) {
			check.Sequential(1)

			t.Run(func(t *T) {
				check.Sequential(2)

				for i := 0; i < 10; i++ {
					t.Run(func(t *T) {
						check.Sequential(3 + i)

						t.Parallel()

						check.Parallel(14, 10)
					})
				}

				check.Sequential(13)
			})

			check.Sequential(15)
		},
		func(t *T) {
			check.Sequential(16)
		},
	})
}

func TestTwoParallelSequences(t *testing.T) {
	check := &ConcurrencyChecker{t: t}
	defer check.Finish(4)

	Run([]func(*T){
		func(t *T) {
			check.Sequential(0)
			t.Parallel()

			t.Run(func(t *T) {
				check.Parallel(2, 2)
			})

			t.Run(func(t *T) {
				check.Parallel(3, 2)
			})
		},
		func(t *T) {
			check.Sequential(1)
			t.Parallel()

			t.Run(func(t *T) {
				check.Parallel(2, 2)
			})

			t.Run(func(t *T) {
				check.Parallel(3, 2)
			})
		},
	})
}
