package tests

import (
	"context"
	"testing"
	"time"

	"rate-limiter/limiter"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	ctx := context.Background()
	// Use the Redis address as a string
	redisAddr := "localhost:6379"
	storage := limiter.NewRedisStorage(redisAddr)
	rateLimiter := limiter.NewLimiter(storage)

	// Test limitation by IP
	ip := "192.168.1.1"
	limit := 5

	for i := 0; i < limit; i++ {
		allowed, err := rateLimiter.Allow(ctx, ip, limit, time.Second*5)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	// Test block after exceeding limit
	allowed, err := rateLimiter.Allow(ctx, ip, limit, time.Second*5)
	assert.NoError(t, err)
	assert.False(t, allowed)
}
