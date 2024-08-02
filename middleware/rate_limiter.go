package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		log.Println("Requisição recebida:", r.RemoteAddr)
		token := r.Header.Get("API_KEY")
		//ipPort := r.RemoteAddr
		//ip, _, err := net.SplitHostPort(ipPort)
		var ip string

		forwardedFor := r.Header.Get("X-Forwarded-For")
		if forwardedFor != "" {
			ip = strings.Split(forwardedFor, ",")[0] // Pega o primeiro IP da lista
		} else {
			// Se não houver cabeçalho X-Forwarded-For, usa o IP de conexão direta
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}

		ctx := context.Background()
		maxRequestsPerSecondToken, _ := strconv.Atoi(os.Getenv("MAX_REQUESTS_PER_SECOND_TOKEN"))
		blockTimeSeconds, _ := strconv.Atoi(os.Getenv("BLOCK_TIME_SECONDS"))

		var allowed bool
		var err error

		if token != "" {
			log.Println("Requisição permitida para: token:", token)
			allowed, err = rateLimiter.Allow(ctx, "token:"+token, maxRequestsPerSecondToken, time.Duration(blockTimeSeconds)*time.Second)
		} else {
			log.Println("Requisição permitida para: ip:", ip)
			maxRequestsPerSecondIP, _ := strconv.Atoi(os.Getenv("MAX_REQUESTS_PER_SECOND_IP"))
			allowed, err = rateLimiter.Allow(ctx, "ip:"+ip, maxRequestsPerSecondIP, time.Duration(blockTimeSeconds)*time.Second)
		}

		if err != nil || !allowed {
			if token != "" {
				log.Println("Requisição bloqueada para: token:", token)
			} else {
				log.Println("Requisição bloqueada para: ip:", ip)
			}
			http.Error(w, "You have reached the maximum number of requests allowed within a certain time frame.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
