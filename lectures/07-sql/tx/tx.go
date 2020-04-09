package tx

import (
	"context"
	"database/sql"
	"log"
)

func Begin(ctx context.Context, db *sql.DB) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `UPDATE users SET name = "Tyador Borlú" WHERE id = 1`); err != nil {
		log.Fatal(err)
	}

	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
