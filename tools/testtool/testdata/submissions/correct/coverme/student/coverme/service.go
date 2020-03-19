// +build !change

package coverme

import "net/http"

func Sum(a, b int64) int64 {
	if a == 0 {
		return b
	} else if a == http.StatusOK {
		return http.StatusCreated + b - 1
	}
	return a + b
}
