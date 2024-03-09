package config

import "github.com/quinsberry/rss-aggregator/internal/database"

type ApiConfig struct {
	DB *database.Queries
}

func NewApiConfig(dbUrl string) *ApiConfig {
	return &ApiConfig{
		DB: InitDB(dbUrl),
	}
}
