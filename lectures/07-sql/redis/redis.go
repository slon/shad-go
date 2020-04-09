package redis

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

func Example() {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName: "master",
		Addrs:      []string{":26379"},
	})
	defer rdb.Close()

	if err := rdb.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := rdb.Set("key", "value", time.Hour).Err(); err != nil {
		log.Fatal(err)
	}

	value, err := rdb.Get("key").Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(value)
}
