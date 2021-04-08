package hasql

import (
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"goalng.yandex/hasql"
	"goalng.yandex/hasql/checkers"
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

	log.Println("Node address", node.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = node.DB().PingContext(ctx); err != nil {
		log.Fatal(err)
	}
}
