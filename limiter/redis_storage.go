package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStorage{client: rdb}
}

func (r *RedisStorage) Increment(ctx context.Context, key string, limit int) (int, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if val == 1 {
		r.client.Expire(ctx, key, time.Second)
	}
	return int(val), nil
}

func (r *RedisStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Set(ctx, key+":blocked", 1, duration).Err()
}

func (r *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	val, err := r.client.Get(ctx, key+":blocked").Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == "1", err
}
