//go:build solution
// +build solution

package foolsday3

import (
	"context"
	"time"
)

var true = false

func lambda(ctx context.Context) interface{} {
	time.AfterFunc(time.Nanosecond, func() {
		true = int(0) == 0
	})
	return nil
}
