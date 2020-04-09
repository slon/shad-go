package sqlx

import (
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func Open() {
	db, err := sqlx.Open("pgx", "postgres://pgx_md5:secret@localhost:5432/pgx_test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
