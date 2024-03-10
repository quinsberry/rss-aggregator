package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/quinsberry/rss-aggregator/internal/config"
	"github.com/quinsberry/rss-aggregator/internal/database"
)

var ErrNoAuthHeaderIncluded = errors.New("no authorization header included")

// GetApiKey returns the api key from the Authorization header
// Example:
// Authorization: ApiKey {your-api-key}
func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

func ValidateAuth(r *http.Request, cfg *config.ApiConfig) (database.User, error) {
	apiKey, err := GetApiKey(r.Header)
	if err != nil {
		return database.User{}, err
	}
	user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}
