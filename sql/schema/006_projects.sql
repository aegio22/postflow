-- +goose Up
ALTER TABLE projects
    ADD CONSTRAINT projects_title_unique UNIQUE (title);

-- +goose Down
ALTER TABLE projects
    DROP CONSTRAINT IF EXISTS projects_title_unique;