package tour1

import (
	"testing"
)

func TestTour(t *testing.T) {
	expected := "tour1 done!"
	if out := Tour(); out != expected {
		t.Errorf("expected %q got %q", expected, out)
	}
}
