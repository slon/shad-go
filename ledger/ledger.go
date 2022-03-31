//go:build !solution

package ledger

import "context"

func New(ctx context.Context, dsn string) (Ledger, error) {
	panic("not implemented")
}
