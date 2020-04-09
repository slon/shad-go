package resources

import (
	"context"
	"database/sql"
	"log"
)

func RowsExhaust(ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM users")
	if err != nil {
		log.Fatal(err)
	}

	if rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}

		log.Println(id, name)
	}
}
