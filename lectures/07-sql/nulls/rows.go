package nulls

import (
	"context"
	"database/sql"
	"log"
)

func Results(ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM users WHERE id = $1", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var s sql.NullString
		if err := rows.Scan(&s); err != nil {
			log.Fatal(err)
		}

		if s.Valid {
			//
		} else {
			//
		}
	}
}
