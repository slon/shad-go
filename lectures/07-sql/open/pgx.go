import (
    "context"

    "github.com/jackc/pgx/v4"
)

conn, err := pgx.Connect(context.Background(), "postgres://pgx_md5:secret@localhost:5432/pgx_test")
