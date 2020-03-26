// +build solution

package big

import "testing"

func TestFunc7(t *testing.T) {
	if got, want := Func7(), 7; got != want {
		t.Errorf("Func7() = %d; want %d", got, want)
	}
}
