// +build solution

package sum

import "testing"

func TestSum(t *testing.T) {
	x, y := 1, 1
	if got, want := Sum(x, y), 2; got != want {
		t.Errorf("Sum(%d, %d) = %d; want %d", x, y, got, want)
	}
}
