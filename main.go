package main

import (
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/joho/godotenv"
	"github.com/quinsberry/rss-aggregator/internal/config"
	"github.com/quinsberry/rss-aggregator/internal/router"
)

func main() {
	godotenv.Load()

	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	cfg := config.NewApiConfig(dbUrl)
	r := router.NewRouter(cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", portStr),
		Handler: r,
	}

	log.Printf("Server is running on port %s", portStr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Server is shutting down")
	}
}
