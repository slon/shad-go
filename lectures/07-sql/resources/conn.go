package resources

import (
	"context"
	"database/sql"
	"log"
)

func ConnExhaust(ctx context.Context, db *sql.DB) {
	c, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_ = c.PingContext(ctx)
}
