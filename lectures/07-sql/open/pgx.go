package open

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func PGXOpen() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://pgx_md5:secret@localhost:5432/pgx_test")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)
}
