package sql

import (
	"context"
	"database/sql/driver"
)

type Driver interface {
	Open(name string) (driver.Conn, error)
}

type QueryerContext interface {
	QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error)
}
