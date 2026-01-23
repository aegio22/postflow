-- name: AddNewProjectUser :one
INSERT INTO users_projects(project_id, user_id, user_status)
VALUES(
    $1,
    $2,
    $3
)
RETURNING *;


-- name: GetUserProjectRelation :one
SELECT * FROM users_projects
WHERE user_id = $1 AND project_id = $2;

-- name: RemoveUserFromProject :exec
DELETE FROM users_projects
WHERE user_id = $1 AND project_id = $2;