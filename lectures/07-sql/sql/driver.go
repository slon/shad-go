type Driver interface {
    Open(name string) (Conn, error)
}

type QueryerContext interface {
    QueryContext(ctx context.Context, query string, args []NamedValue) (Rows, error)
}
