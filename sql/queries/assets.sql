-- name: CreateAsset :one
INSERT INTO assets (
    project_id,
    name,
    storage_path,
    tags,
    created_by
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;


-- name: UpdateAssetStoragePath :exec
UPDATE assets
SET storage_path = $2, updated_at = NOW()
WHERE id = $1;


-- name: GetAssetByName :one
SELECT * FROM assets
WHERE name = $1 AND project_id= $2;

-- name: GetAssetsByProjectName :many
SELECT a.* 
from assets a
JOIN projects p
    on p.id = a.project_id
WHERE p.title = $1;