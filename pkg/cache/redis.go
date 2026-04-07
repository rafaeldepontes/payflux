package cache

import (
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func GetCache() *redis.Client {
	once.Do(func() {
		idx, err := strconv.Atoi(os.Getenv("REDIS_DB_IDX"))
		if err != nil {
			idx = 0
		}

		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       idx,
		})
	})

	return rdb
}

func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
