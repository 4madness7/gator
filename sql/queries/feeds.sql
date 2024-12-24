-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedsWithUser :many
SELECT users.name AS user_name, feeds.name AS feed_name, feeds.url FROM feeds
INNER JOIN users
ON feeds.user_id = users.id;

-- name: GetFeedWithURL :one
SELECT * FROM feeds WHERE url = $1 LIMIT 1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2, last_fetched_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST, updated_at
LIMIT 1;
