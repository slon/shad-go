package cancelation

import (
	"context"
	"time"
)

func SimpleCancelation() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	if err := doSlowJob(ctx); err != nil {
		panic(err)
	}
}

// OMIT

func SimpleTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := doSlowJob(ctx); err != nil {
		panic(err)
	}
}

// OMIT

func doSlowJob(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// perform a portion of slow job
			time.Sleep(1 * time.Second)
		}
	}
}

// OMIT
