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