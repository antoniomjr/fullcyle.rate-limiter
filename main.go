package main

import (
	"log"
	"net/http"
	"os"

	"rate-limiter/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := mux.NewRouter()

	r.Use(middleware.RateLimiterMiddleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port", port)
	http.ListenAndServe(":"+port, r)
}
