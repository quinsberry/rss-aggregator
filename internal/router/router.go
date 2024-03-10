package router

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/quinsberry/rss-aggregator/internal/config"
	"github.com/quinsberry/rss-aggregator/internal/handlers"
)

func NewRouter(cfg *config.ApiConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/users", handlers.HandlerCreateUser(cfg))
		r.Get("/users", handlers.HandlerGetUser(cfg))
		r.Post("/feeds", handlers.HandlerCreateFeed(cfg))
		r.Get("/feeds", handlers.HandlerGetFeeds(cfg))
	})
	return r
}
