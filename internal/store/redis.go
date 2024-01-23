package store

import (
	"os"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Rdb *redis.Client
}

func NewRedis() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Redis{
		Rdb: rdb,
	}
}
