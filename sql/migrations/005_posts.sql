-- +goose Up 
CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  url VARCHAR(255) UNIQUE NOT NULL,
  description TEXT,
  published_at TIMESTAMP ,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS posts;
