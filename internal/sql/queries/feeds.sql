-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, user_id, url, name)
VALUES ($1, $2, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedById :one
SELECT * FROM feeds WHERE id = $1;

-- name: DeleteFeed :one
DELETE FROM feeds WHERE id = $1
RETURNING *;