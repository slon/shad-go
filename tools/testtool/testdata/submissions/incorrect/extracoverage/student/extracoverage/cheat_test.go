package extracoverage

import "testing"

// Gives package coverage 75%
func TestLongFunc(t *testing.T) {
	if got, want := LongFunc(), 7; got != want {
		t.Errorf("LongFunc() = %d; want %d", got, want)
	}
}
