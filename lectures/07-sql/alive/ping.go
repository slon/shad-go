package alive

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func IsItAliveQuestionMark(ctx context.Context) {
	db, err := sql.Open("pgx", "postgres://pgx_md5:secret@localhost:5432/pgx_test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
}
