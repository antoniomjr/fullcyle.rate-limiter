package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"

	"rate-limiter/middleware"

	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := mux.NewRouter()

	r.Use(middleware.RateLimiterMiddleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Não foi possível conectar ao Redis: %v", err)
	}
	fmt.Println("Conexão com o Redis estabelecida com sucesso!")

	log.Println("Server starting on port", port)
	http.ListenAndServe(":"+port, r)
}
