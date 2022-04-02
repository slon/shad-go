//go:build !solution

package dao

import "context"

func CreateDao(ctx context.Context, dsn string) (Dao, error) {
	panic("not implemented")
}
