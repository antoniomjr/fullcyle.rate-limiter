package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"rate-limiter/limiter"

	"github.com/joho/godotenv"
)

var rateLimiter *limiter.Limiter

func init() {
	godotenv.Load()

	redisAddr := os.Getenv("REDIS_ADDR")
	rateLimiter = limiter.NewLimiter(limiter.NewRedisStorage(redisAddr))
}

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")

		ctx := context.Background()
		maxRequestsPerSecondToken, _ := strconv.Atoi(os.Getenv("MAX_REQUESTS_PER_SECOND_TOKEN"))
		blockTimeSeconds, _ := strconv.Atoi(os.Getenv("BLOCK_TIME_SECONDS"))

		var allowed bool
		var err error

		if token != "" {
			allowed, err = rateLimiter.Allow(ctx, "token:"+token, maxRequestsPerSecondToken, time.Duration(blockTimeSeconds)*time.Second)
		} else {
			ip := r.RemoteAddr
			maxRequestsPerSecondIP, _ := strconv.Atoi(os.Getenv("MAX_REQUESTS_PER_SECOND_IP"))
			allowed, err = rateLimiter.Allow(ctx, "ip:"+ip, maxRequestsPerSecondIP, time.Duration(blockTimeSeconds)*time.Second)
		}

		if err != nil || !allowed {
			http.Error(w, "You have reached the maximum number of requests allowed within a certain time frame.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
