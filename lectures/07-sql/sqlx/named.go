package sqlx

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

func Insert(ctx context.Context, db *sqlx.DB) {
	_, err := db.NamedExecContext(
		ctx,
		"INSERT INTO users(name) VALUES(:name)",
		map[string]interface{}{
			"name": "Jukka Sarasti",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
