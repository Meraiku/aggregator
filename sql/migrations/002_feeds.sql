-- +goose Up

CREATE TABLE IF NOT EXISTS feeds(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_name ON feeds USING hash (name);

-- +goose Down

DROP TABLE IF EXISTS feeds;