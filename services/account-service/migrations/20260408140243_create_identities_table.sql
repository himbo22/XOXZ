-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS identities (
  id UUID PRIMARY KEY,
  user_id UUID,
  provider VARCHAR(50), -- google, facebook, zalo, stripchat, ...
  provider_user_id VARCHAR(255), -- ID from third party
  provider_data JSONB,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_identities_provider_uid ON identities(provider, provider_user_id);
CREATE INDEX idx_identities_user_id ON identities(user_id) WHERE deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS identities;
-- +goose StatementEnd

