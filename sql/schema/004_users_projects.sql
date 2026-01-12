-- +goose Up
CREATE TABLE users_projects(
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_status TEXT NOT NULL
);

CREATE INDEX idx_projects_user_status ON users_projects(user_status);
CREATE INDEX idx_projects_user_id ON users_projects(user_id);
CREATE INDEX idx_projects_project_id ON users_projects(project_id);

-- +goose Down
DROP INDEX idx_projects_project_id ON users_projects(project_id) IF EXISTS;
DROP INDEXidx_projects_user_id ON users_projects(user_id) IF EXISTS;
DROP INDEX idx_projects_user_status ON users_projects(user_status) IF EXISTS;
DROP TABLE users_projects IF EXISTS;