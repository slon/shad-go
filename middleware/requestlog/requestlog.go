//go:build !solution

package requestlog

import (
	"net/http"

	"go.uber.org/zap"
)

func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	panic("not implemented")
}
