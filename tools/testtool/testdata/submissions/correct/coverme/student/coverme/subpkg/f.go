// +build !change

package subpkg

func AddOne(n int) int {
	if n == 0 {
		return 1
	} else if n == 1 {
		return 2
	}
	return n + 1
}
