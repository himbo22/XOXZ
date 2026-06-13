-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS artist_groups (
    id UUID PRIMARY KEY,
    display_name VARCHAR(120),
    bio TEXT,
    avatar_url TEXT,
    banner_url TEXT,
    verified BOOLEAN DEFAULT FALSE,
    status VARCHAR(32) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS artist_groups;
-- +goose StatementEnd
