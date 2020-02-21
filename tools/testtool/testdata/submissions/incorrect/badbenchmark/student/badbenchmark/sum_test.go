package badbenchmark

import (
	"fmt"
	"math"
	"testing"
)

type testCase struct {
	a, b, sum int64
}

func BenchmarkSum(b *testing.B) {
	for i, input := range []testCase{
		{a: 2, b: 2, sum: 4},
		{a: 2, b: -2, sum: 0},
		{a: math.MaxInt64, b: 1, sum: math.MinInt64},
	} {
		b.Run(fmt.Sprint(i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Sum(input.a, input.b)
			}
		})
	}
}
