-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS artists (
    id UUID PRIMARY KEY,
    user_id UUID,
    group_id UUID,
    stage_name VARCHAR(120),
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

CREATE INDEX IF NOT EXISTS idx_artists_user_id ON artists(user_id);
CREATE INDEX IF NOT EXISTS idx_artists_deleted_at ON artists(deleted_at);
CREATE INDEX IF NOT EXISTS idx_artists_stage_name ON artists(stage_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS artists;
-- +goose StatementEnd
