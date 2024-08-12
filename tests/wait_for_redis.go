// tests/wait_for_redis.go
package tests

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func WaitForRedis() {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	for {
		_, err := client.Ping(context.Background()).Result()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
}
