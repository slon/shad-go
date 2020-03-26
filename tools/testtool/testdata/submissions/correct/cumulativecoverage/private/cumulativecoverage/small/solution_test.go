// +build solution

package small

import "testing"

func TestFunc3(t *testing.T) {
	if got, want := Func3(), 3; got != want {
		t.Errorf("Func3() = %d; want %d", got, want)
	}
}
