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


-- name: GetProjectByTitle :one
SELECT * FROM projects
WHERE title = $1;

-- name: DeleteProjectByTitle :exec
DELETE FROM projects
WHERE title = $1;

-- name: GetProjectsForUser :many
SELECT p.*
FROM projects p
JOIN users_projects up
  ON up.project_id = p.id
WHERE up.user_id = $1;