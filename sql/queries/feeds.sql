-- name: CreateFeed :one

INSERT INTO feeds (id, name, url, user_id, created_at, updated_at) 
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, url, user_id;

-- name: GetAllFeeds :many

SELECT  feeds.name AS feed_name, 
        feeds.url, 
        users.name AS user_name
FROM feeds
INNER JOIN users
ON feeds.user_id=users.id;

-- name: GetFeedIDByURL :one

SELECT id FROM feeds
WHERE url = $1;

-- name: MarkFetched :exec

UPDATE feeds
SET last_fetched_at = $1,
    updated_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one

SELECT id, url
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
