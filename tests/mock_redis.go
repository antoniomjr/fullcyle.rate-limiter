// tests/redis_mock.go
package tests

import (
	"context"
	"sync"
	"time"
)

type MockRedisStorage struct {
	mu    sync.Mutex
	store map[string]time.Time
	count map[string]int
}

func NewMockRedisStorage() *MockRedisStorage {
	return &MockRedisStorage{
		store: make(map[string]time.Time),
		count: make(map[string]int),
	}
}

func (m *MockRedisStorage) Set(ctx context.Context, key string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = time.Now().Add(expiration)
	return nil
}

func (m *MockRedisStorage) Get(ctx context.Context, key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	expiration, exists := m.store[key]
	if !exists {
		return false, nil
	}
	if time.Now().After(expiration) {
		delete(m.store, key)
		return false, nil
	}
	return true, nil
}

func (m *MockRedisStorage) Increment(ctx context.Context, key string, limit int) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.count[key]++
	return m.count[key], nil
}

func (m *MockRedisStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	return m.Set(ctx, key, duration)
}

func (m *MockRedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	return m.Get(ctx, key)
}
