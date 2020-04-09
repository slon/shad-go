package prepare

import (
	"context"
	"database/sql"
	"log"
)

func Prepare(ctx context.Context, db *sql.DB) {
	stmt, err := db.PrepareContext(ctx, "SELECT name FROM users WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 1; ; i++ {
		var name string
		if err = stmt.QueryRowContext(ctx, i).Scan(&name); err != nil {
			log.Fatal(err)
		}

		log.Println(name)
	}
}
