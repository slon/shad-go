package tour0

import (
	"testing"
)

func TestTour(t *testing.T) {
	expected := "tour0 done!"
	if out := Tour(); out != expected {
		t.Errorf("expected %q got %q", expected, out)
	}
}
