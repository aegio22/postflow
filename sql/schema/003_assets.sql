-- +goose Up
CREATE TABLE assets(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    asset_type TEXT NOT NULL,
    storage_path TEXT NOT NULL,
    tags TEXT[],
    current_version_number INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft',
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_assets_project ON assets(project_id);
CREATE INDEX idx_assets_created_by ON assets(created_by);
CREATE INDEX idx_assets_type ON assets(asset_type);
CREATE INDEX idx_assets_status ON assets(status);

-- Composite index for common query: "Show all audio mixes in this project"
CREATE INDEX idx_assets_project_type ON assets(project_id, asset_type);

-- +goose Down
DROP INDEX IF EXISTS idx_assets_project_type;
DROP INDEX IF EXISTS idx_assets_status;
DROP INDEX IF EXISTS idx_assets_type;
DROP INDEX IF EXISTS idx_assets_created_by;
DROP INDEX IF EXISTS idx_assets_project;
DROP TABLE IF EXISTS assets;