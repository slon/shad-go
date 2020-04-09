package sqlx

import (
    "context"
    "database/sql"

    "github.com/jmoiron/sqlx"
)

type QueryerContext interface {
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
    QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type ExecerContext interface {
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type ExtContext interface {
    QueryerContext
    ExecerContext
}
