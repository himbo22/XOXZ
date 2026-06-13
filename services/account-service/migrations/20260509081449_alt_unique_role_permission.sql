-- +goose Up
-- +goose StatementBegin
ALTER TABLE roles
    ADD CONSTRAINT idx_roles_name_unique UNIQUE (name);

ALTER TABLE permissions
    ADD CONSTRAINT idx_permissions_code_unique UNIQUE (code);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE roles
DROP CONSTRAINT IF EXISTS idx_roles_name_unique;

ALTER TABLE permissions
DROP CONSTRAINT IF EXISTS idx_permissions_code_unique;
-- +goose StatementEnd