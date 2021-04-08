package resources

import (
	"context"
	"database/sql"
	"log"
)

func NoContext(ctx context.Context, db *sql.DB) {
	// У Conn() нет версии без контекста
	c, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Потенциально вечный Ping
	_ = c.Ping()
}
