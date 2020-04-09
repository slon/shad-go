package query

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

func QueryRow(ctx context.Context, db *sql.DB) {
	var id int
	var name string
	err := db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", 1).Scan(&id, &name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("nothing found")
			return
		}

		log.Fatal(err)
	}

	log.Println(name)
}
