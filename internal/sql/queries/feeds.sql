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

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds ORDER BY last_fetched_at DESC LIMIT $1;

-- name: MarkFeedFetched :one
UPDATE feeds SET updated_at = $1, last_fetched_at = $1 WHERE id = $2
RETURNING *;