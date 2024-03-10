package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/quinsberry/rss-aggregator/internal/auth"
	"github.com/quinsberry/rss-aggregator/internal/config"
	"github.com/quinsberry/rss-aggregator/internal/database"
	"github.com/quinsberry/rss-aggregator/internal/utils"
)

func HandlerCreateFeed(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
		type response struct {
			Id        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			UserId    uuid.UUID `json:"user_id"`
			Name      string    `json:"name"`
			URL       string    `json:"url"`
		}

		params := parameters{}
		err := utils.DecodeJSONBody(w, r, &params)
		if err != nil {
			var mr *utils.MalformedRequest
			if errors.As(err, &mr) {
				utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
			} else {
				log.Print(err.Error())
				utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
			return
		}

		user, err := auth.ValidateAuth(r, cfg)
		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
			return
		}

		feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UserID:    user.ID,
			Name:      params.Name,
			Url:       params.URL,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating feed: %v", err))
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, response{
			Id:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			UserId:    feed.UserID,
			Name:      feed.Name,
			URL:       feed.Url,
		})
	}
}

func HandlerGetFeeds(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Id        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			UserId    uuid.UUID `json:"user_id"`
			Name      string    `json:"name"`
			URL       string    `json:"url"`
		}

		feeds, err := cfg.DB.GetAllFeeds(r.Context())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting feeds: %v", err))
			return
		}

		responseFeeds := make([]response, len(feeds))
		for _, feed := range feeds {
			responseFeeds = append(responseFeeds, response{
				Id:        feed.ID,
				CreatedAt: feed.CreatedAt,
				UpdatedAt: feed.UpdatedAt,
				UserId:    feed.UserID,
				Name:      feed.Name,
				URL:       feed.Url,
			})
		}

		utils.RespondWithJSON(w, http.StatusCreated, responseFeeds)
	}
}
