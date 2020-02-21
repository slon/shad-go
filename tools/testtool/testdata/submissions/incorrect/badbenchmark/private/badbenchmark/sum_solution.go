// +build solution

package badbenchmark

import "time"

func Sum(a, b int64) int64 {
	time.Sleep(time.Millisecond)
	return a + b
}
