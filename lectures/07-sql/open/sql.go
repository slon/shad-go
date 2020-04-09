package open

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func SQLOpen() {
	db, err := sql.Open("pgx", "postgres://pgx_md5:secret@localhost:5432/pgx_test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
