-- +goose Up
ALTER TABLE assets
    ADD CONSTRAINT assets_project_name_unique UNIQUE (project_id, name);

-- +goose Down
ALTER TABLE assets
    DROP CONSTRAINT IF EXISTS assets_project_name_unique;