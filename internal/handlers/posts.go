package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/quinsberry/rss-aggregator/internal/auth"
	"github.com/quinsberry/rss-aggregator/internal/config"
	"github.com/quinsberry/rss-aggregator/internal/database"
	"github.com/quinsberry/rss-aggregator/internal/utils"
)

func HandlerGetPostsByUser(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Id          uuid.UUID `json:"id"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
			Title       string    `json:"title"`
			Url         string    `json:"url"`
			Description string    `json:"description"`
			PublishedAt time.Time `json:"published_at"`
			FeedId      uuid.UUID `json:"feed_id"`
		}

		limitQry := r.URL.Query().Get("limit")
		limit := 10

		log.Println(int32(limit))

		var err error
		if limitQry != "" {
			limit, err = strconv.Atoi(limitQry)
			if err != nil {
				log.Printf("Error formating limit query defaulting to 10: %v", err)
				limit = 10
			}
		}

		user, err := auth.ValidateAuth(r, cfg)
		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
		}

		posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprint(err))
			return
		}

		responsePosts := make([]response, len(posts))
		for i, post := range posts {
			responsePosts[i] = response{
				Id:          post.ID,
				CreatedAt:   post.CreatedAt,
				UpdatedAt:   post.UpdatedAt,
				Title:       post.Title,
				Url:         post.Url,
				FeedId:      post.FeedID,
				PublishedAt: post.PublishedAt,
				Description: post.Description.String,
			}
		}

		utils.RespondWithJSON(w, http.StatusOK, responsePosts)
	}
}
