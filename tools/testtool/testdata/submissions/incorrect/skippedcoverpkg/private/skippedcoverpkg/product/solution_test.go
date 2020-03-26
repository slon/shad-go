// +build solution

package product

import "testing"

func TestProduct(t *testing.T) {
	x, y := 3, 3
	if got, want := Product(x, y), 9; got != want {
		t.Errorf("Product(%d, %d) = %d; want %d", x, y, got, want)
	}
}
