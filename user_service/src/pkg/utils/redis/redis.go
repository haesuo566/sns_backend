package redis

import (
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Client = redis.Client

var once sync.Once
var instance *redis.Client

func New() *redis.Client {
	once.Do(func() {
		instance = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})
	return instance
}
