package redis

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Util interface {
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	Get(context.Context, string) *redis.StringCmd
	Del(context.Context, ...string) *redis.IntCmd
	SetNX(context.Context, string, interface{}, time.Duration) *redis.BoolCmd
}

var once sync.Once
var instance Util

func New() Util {
	once.Do(func() {
		instance = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})
	return instance
}
