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

func HandlerGetUser(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Id        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Name      string    `json:"name"`
			APIKey    string    `json:"api_key"`
		}

		user, err := auth.ValidateAuth(r, cfg)
		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
		}

		utils.RespondWithJSON(w, 200, response{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
			APIKey:    user.ApiKey,
		})
	}
}

func HandlerCreateUser(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Name string `json:"name"`
		}
		type response struct {
			Id        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Name      string    `json:"name"`
			ApiKey    string    `json:"api_key"`
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

		user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create user: %v", err))
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, response{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
			ApiKey:    user.ApiKey,
		})
	}
}
