-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;


-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND refresh_tokens.revoked_at IS NULL
AND refresh_tokens.expires_at > NOW();

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1;

-- name: DeleteAllRefreshTokens :exec
DELETE FROM refresh_tokens;


-- name: GetTokenFromUserID :one
SELECT * FROM refresh_tokens 
WHERE user_id = $1;