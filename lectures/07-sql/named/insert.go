package named

import (
	"context"
	"database/sql"
	"log"
)

func Insert(ctx context.Context, db *sql.DB) {
	_, err := db.ExecContext(
		ctx,
		"INSERT INTO users(name) VALUES(@name)",
		sql.Named("name", "Amos Burton"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
