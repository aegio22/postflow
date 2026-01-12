-- name: CreateUser :one
INSERT INTO users (id, username, created_at, updated_at, email, hashed_password)
VALUES(
    gen_random_uuid(),
    $1,
    NOW(),
    NOW(),
    $2, 
    $3
)
RETURNING *;


-- name: DeleteAllUsers :exec
DELETE FROM users;


