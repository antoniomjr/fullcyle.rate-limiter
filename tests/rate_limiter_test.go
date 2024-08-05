// tests/rate_limiter_test.go
package tests

import (
	"context"
	"rate-limiter/limiter"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiterIp(t *testing.T) {
	ctx := context.Background()
	storage := NewMockRedisStorage()
	rateLimiter := limiter.NewLimiter(storage)

	ip := "192.168.1.1"
	limit := 5

	for i := 0; i < limit; i++ {
		allowed, err := rateLimiter.Allow(ctx, ip, limit, time.Second*5)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rateLimiter.Allow(ctx, ip, limit, time.Second*5)
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestRateLimiterToken(t *testing.T) {
	ctx := context.Background()
	storage := NewMockRedisStorage()
	rateLimiter := limiter.NewLimiter(storage)

	token := "your_token"
	limit := 5

	for i := 0; i < limit; i++ {
		allowed, err := rateLimiter.Allow(ctx, token, limit, time.Second*5)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rateLimiter.Allow(ctx, token, limit, time.Second*5)
	assert.NoError(t, err)
	assert.False(t, allowed)
}
