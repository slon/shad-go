// +build !solution

package badbenchmark

import "time"

func Sum(a, b int64) int64 {
	time.Sleep(time.Millisecond * 10)
	return a + b
}
