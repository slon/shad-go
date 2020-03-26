// +build !change

package extracoverage

func Sum(a, b int64) int64 {
	if a == 0 {
		return b
	}
	return a + b
}
