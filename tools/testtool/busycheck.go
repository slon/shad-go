package testtool

import (
	"bytes"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func VerifyNoBusyGoroutines(t *testing.T) {
	time.Sleep(time.Millisecond * 100)

	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond)

		var stacks []byte
		for n := 1 << 20; true; n *= 2 {
			stacks = make([]byte, n)
			m := runtime.Stack(stacks, true)

			if m < n {
				stacks = stacks[:m]
				break
			}
		}

		busy := bytes.Count(stacks, []byte("[running]"))
		busy += bytes.Count(stacks, []byte("[runnable]"))
		busy += bytes.Count(stacks, []byte("[sleep]"))

		if !assert.Less(t, busy, 2) {
			_, _ = os.Stderr.Write(stacks)
			break
		}
	}

}
