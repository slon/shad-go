//go:build !solution

package rsem

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Semaphore struct {
}

func NewSemaphore(rdb redis.UniversalClient) *Semaphore {
	panic("not implemented")
}

// Acquire semaphore associated with key. No more than limit processes can hold semaphore at the same time.
func (s *Semaphore) Acquire(
	ctx context.Context,
	key string,
	limit int,
) (release func() error, err error) {
	panic("not implemented")
}
