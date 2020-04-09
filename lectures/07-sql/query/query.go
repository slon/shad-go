package query

import (
	"context"
	"database/sql"
	"log"
)

func Query(ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM users WHERE id = $1", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}

		log.Println(id, name)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
