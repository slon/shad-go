// +build !solution

package newdependency

import "rsc.io/quote/v3"

func Sum(a, b int64) int64 {
	quote.HelloV3()
	return a + b
}
