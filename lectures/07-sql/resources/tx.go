package resources

import (
	"context"
	"database/sql"
	"log"
)

func TxExhaust(ctx context.Context, db *sql.DB) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err = tx.ExecContext(ctx, `UPDATE users SET name = "Surl/Tesh-echer" WHERE id = 1`); err != nil {
		log.Println(err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
	}
}

func TxDeadlock(ctx context.Context, db *sql.DB) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	_, _ = db.QueryContext(ctx, "SELECT id, name FROM users")
	_, _ = db.QueryContext(ctx, "SELECT id, name FROM users")
}
