import (
    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v4/stdlib"
)

db, err := sqlx.Open("pgx", "postgres://pgx_md5:secret@localhost:5432/pgx_test")
if err != nil {
    return err
}
defer db.Close()
