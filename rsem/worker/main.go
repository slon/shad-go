//go:build !change

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"gitlab.com/slon/shad-go/rsem"
)

func do() error {
	rdb := redis.NewClient(&redis.Options{Addr: os.Args[1]})
	defer func() { _ = rdb.Close() }()
	sem := rsem.NewSemaphore(rdb)
	ctx := context.Background()

	_, err := sem.Acquire(ctx, "dead", 2)
	if err != nil {
		return err
	}

	select {}
}

func main() {
	if err := do(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
