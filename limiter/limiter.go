package limiter

import (
	"context"
	"time"
)

type Limiter struct {
	storage Storage
}

type Storage interface {
	Increment(ctx context.Context, key string, limit int) (int, error)
	Block(ctx context.Context, key string, duration time.Duration) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}

func NewLimiter(storage Storage) *Limiter {
	return &Limiter{storage: storage}
}

func (l *Limiter) Allow(ctx context.Context, key string, limit int, duration time.Duration) (bool, error) {
	if blocked, err := l.storage.IsBlocked(ctx, key); err != nil || blocked {
		return false, err
	}
	count, err := l.storage.Increment(ctx, key, limit)
	if err != nil {
		return false, err
	}
	if count > limit {
		l.storage.Block(ctx, key, duration)
		return false, nil
	}
	return true, nil
}
