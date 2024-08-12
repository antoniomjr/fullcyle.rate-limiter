package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"rate-limiter/middleware"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLimiter struct {
	mock.Mock
}

func (m *MockLimiter) Allow(ctx context.Context, key string, limit int, blockTime time.Duration) (bool, error) {
	args := m.Called(ctx, key, limit, blockTime)
	return args.Bool(0), args.Error(1)
}

func (m *MockLimiter) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	args := m.Called(ctx, key, expiration)
	return args.Int(0), args.Error(1)
}

func (m *MockLimiter) Block(ctx context.Context, key string, duration time.Duration) error {
	args := m.Called(ctx, key, duration)
	return args.Error(0)
}

func (m *MockLimiter) IsBlocked(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func setupMockLimiter() *MockLimiter {
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("MAX_REQUESTS_PER_SECOND_TOKEN", "5")
	os.Setenv("MAX_REQUESTS_PER_SECOND_IP", "5")
	os.Setenv("BLOCK_TIME_SECONDS", "10")

	mockLimiter := new(MockLimiter)
	return mockLimiter
}

func TestAllowRequestWithToken(t *testing.T) {
	mockLimiter := setupMockLimiter()
	mockLimiter.On("Allow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	mockLimiter.On("Increment", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	mockLimiter.On("Block", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockLimiter.On("IsBlocked", mock.Anything, mock.Anything).Return(false, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "test_token")

	rr := httptest.NewRecorder()
	handler := middleware.RateLimiterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestBlockRequestWithToken(t *testing.T) {
	mockLimiter := setupMockLimiter()
	mockLimiter.On("Allow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	mockLimiter.On("Increment", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	mockLimiter.On("Block", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockLimiter.On("IsBlocked", mock.Anything, mock.Anything).Return(true, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "test_token")

	rr := httptest.NewRecorder()
	handler := middleware.RateLimiterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 10; i++ {
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	require.Equal(t, http.StatusTooManyRequests, rr.Code)
}

func TestAllowRequestWithIP(t *testing.T) {
	mockLimiter := setupMockLimiter()
	mockLimiter.On("Allow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	mockLimiter.On("Increment", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	mockLimiter.On("Block", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockLimiter.On("IsBlocked", mock.Anything, mock.Anything).Return(false, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"

	rr := httptest.NewRecorder()
	handler := middleware.RateLimiterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestBlockRequestWithIP(t *testing.T) {
	mockLimiter := setupMockLimiter()
	mockLimiter.On("Allow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	mockLimiter.On("Increment", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	mockLimiter.On("Block", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockLimiter.On("IsBlocked", mock.Anything, mock.Anything).Return(true, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"

	rr := httptest.NewRecorder()
	handler := middleware.RateLimiterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 10; i++ {
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	require.Equal(t, http.StatusTooManyRequests, rr.Code)
	require.Equal(t, "You have reached the maximum number of requests allowed within a certain time frame.", rr.Body.String())
}
