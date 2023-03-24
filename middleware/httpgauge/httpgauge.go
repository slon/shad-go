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

func (g *Gauge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	panic("not implemented")
}
