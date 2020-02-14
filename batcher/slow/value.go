// +build !change

package slow

import (
	"sync"
	"sync/atomic"
	"time"
)

type Value struct {
	mu          sync.Mutex
	value       interface{}
	readRunning int32
}

func (s *Value) Load() interface{} {
	if atomic.SwapInt32(&s.readRunning, 1) == 1 {
		panic("another load is running")
	}
	defer atomic.StoreInt32(&s.readRunning, 0)

	s.mu.Lock()
	value := s.value
	s.mu.Unlock()

	time.Sleep(time.Millisecond)
	return value
}

func (s *Value) Store(v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.value = v
}
