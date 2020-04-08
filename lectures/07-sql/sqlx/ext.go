type QueryerContext interface {
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryxContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
    QueryRowxContext(ctx context.Context, query string, args ...interface{}) *Row
}

type ExecerContext interface {
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type ExtContext interface {
    QueryerContext
    ExecerContext
}
