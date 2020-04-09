package nulls

import (
	"context"
	"database/sql"
	"log"
)

func Insert(ctx context.Context, db *sql.DB, name interface{}) {
	_, err := db.ExecContext(
		ctx,
		"INSERT INTO users(name) VALUES(@name)",
		sql.Named("name", name),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func DoStuff(ctx context.Context, db *sql.DB) {
	// Nulls
	Insert(ctx, db, nil)
	Insert(ctx, db, sql.NullString{})

	// Values
	Insert(ctx, db, "The Shrike")
	Insert(ctx, db, sql.NullString{String: "The Shrike", Valid: true})
}
