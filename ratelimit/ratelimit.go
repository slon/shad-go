// +build !solution

package ratelimit

import (
	"context"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
}

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	panic("not implemented")
}

func (l *Limiter) Acquire(ctx context.Context) error {
	panic("not implemented")
}
