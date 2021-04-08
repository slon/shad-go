package redis

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func Example(ctx context.Context) {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName: "master",
		Addrs:      []string{":26379"},
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	if err := rdb.Set(ctx, "key", "value", time.Hour).Err(); err != nil {
		log.Fatal(err)
	}

	value, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(value)
}
