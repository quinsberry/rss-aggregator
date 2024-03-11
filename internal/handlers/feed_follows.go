package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/quinsberry/rss-aggregator/internal/auth"
	"github.com/quinsberry/rss-aggregator/internal/config"
	"github.com/quinsberry/rss-aggregator/internal/database"
	"github.com/quinsberry/rss-aggregator/internal/utils"
)

func HandlerCreateFeedFollow(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			FeedId uuid.UUID `json:"feed_id"`
		}
		type response struct {
			Id        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			UserId    uuid.UUID `json:"user_id"`
			FeedId    uuid.UUID `json:"feed_id"`
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

		feed, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    params.FeedId,
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
			FeedId:    feed.FeedID,
		})
	}
}

func HandlerGetUserFeedFollows(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Id        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			UserId    uuid.UUID `json:"user_id"`
			FeedId    uuid.UUID `json:"feed_id"`
		}

		user, err := auth.ValidateAuth(r, cfg)
		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
			return
		}

		feedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating feed: %v", err))
			return
		}

		responseFeedFollows := make([]response, len(feedFollows))
		for i, feedFollow := range feedFollows {
			responseFeedFollows[i] = response{
				Id:        feedFollow.ID,
				CreatedAt: feedFollow.CreatedAt,
				UpdatedAt: feedFollow.UpdatedAt,
				UserId:    feedFollow.UserID,
				FeedId:    feedFollow.FeedID,
			}
		}

		utils.RespondWithJSON(w, http.StatusOK, responseFeedFollows)
	}
}

func HandlerDeleteFeedFollow(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		feedFollowIDStr := chi.URLParam(r, "id")
		feedFollowID, err := uuid.Parse(feedFollowIDStr)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse feed follow id: %v", err))
			return
		}

		user, err := auth.ValidateAuth(r, cfg)
		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
			return
		}

		err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
			ID:     feedFollowID,
			UserID: user.ID,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not delete feed follow: %v", err))
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, struct{}{})
	}
}
