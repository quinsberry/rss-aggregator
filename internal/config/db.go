package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/quinsberry/rss-aggregator/internal/database"
)

func InitDB(dbUrl string) *database.Queries {
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
	return database.New(conn)

}
