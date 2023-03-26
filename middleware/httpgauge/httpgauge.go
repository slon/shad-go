//go:build !solution

package httpgauge

import "net/http"

type Gauge struct{}

func New() *Gauge {
	panic("not implemented")
}

func (g *Gauge) Snapshot() map[string]int {
	panic("not implemented")
}

// ServeHTTP returns accumulated statistics in text format ordered by pattern.
//
// For example:
//
//	/a 10
//	/b 5
//	/c/{id} 7
func (g *Gauge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	panic("not implemented")
}
