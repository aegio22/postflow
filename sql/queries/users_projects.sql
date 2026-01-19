-- name: AddNewProjectUser :one
INSERT INTO users_projects(project_id, user_id, user_status)
VALUES(
    $1,
    $2,
    $3
)
RETURNING *;


-- name: GetUserFromUsersProjects :one
SELECT * FROM users_projects
WHERE user_id = $1;