//go:build !solution

package shopfront

import "github.com/go-redis/redis/v8"

func New(rdb *redis.Client) Counters {
	panic("not implemented")
}
