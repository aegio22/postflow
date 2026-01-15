-- name: CreateProject :one
INSERT INTO projects(id, title, description, created_by, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    NULLIF($2, ''),
    $3,
    NOW(),
    NOW()
)
RETURNING *;
