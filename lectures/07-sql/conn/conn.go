package query

import (
	"context"
	"database/sql"
	"log"
)

func Conn(ctx context.Context, db *sql.DB) {
	c, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	_ = c.PingContext(ctx)
}
