package hasql

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.yandex/hasql"
	"golang.yandex/hasql/checkers"
)

func Open() {
	dbFoo, _ := sql.Open("pgx", "host=foo")
	dbBar, _ := sql.Open("pgx", "host=bar")
	cluster, err := hasql.NewCluster(
		[]hasql.Node{hasql.NewNode("foo", dbFoo), hasql.NewNode("bar", dbBar)},
		checkers.PostgreSQL,
	)
	if err != nil {
		log.Fatal(err)
	}

	node := cluster.Primary()
	if err == nil {
		log.Fatal(err)
	}

	log.Println("Node address", node.Addr())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = node.DB().PingContext(ctx); err != nil {
		log.Fatal(err)
	}
}
