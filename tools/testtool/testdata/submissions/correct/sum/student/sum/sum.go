// +build !solution

package sum

import "gitlab.com/slon/shad-go/sum/pkg"

func Sum(a, b int64) int64 {
	pkg.F()
	return a + b
}
