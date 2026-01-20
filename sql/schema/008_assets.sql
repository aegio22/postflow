-- +goose Up
ALTER TABLE assets
    ADD CONSTRAINT assets_name_unique UNIQUE (name);

-- +goose Down
ALTER TABLE assets
    DROP CONSTRAINT IF EXISTS assets_name_unique;