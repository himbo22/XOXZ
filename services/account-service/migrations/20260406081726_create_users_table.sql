-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(36),
  email VARCHAR(255),
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  phone VARCHAR(50),
  avatar_url VARCHAR(255),
  bio VARCHAR(255),
  last_login_at TIMESTAMP,
  status VARCHAR(50), -- ACTIVE, LOCKED
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- ID is UUIDv7 , so we don't need to add index to created_at, because we can use ID as created_at
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_unique
ON users(username)
WHERE username IS NOT NULL AND deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
