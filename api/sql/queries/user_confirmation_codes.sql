-- name: CreateUserLoginCode :one
INSERT INTO user_login_codes (email, code)
VALUES ($1, $2) RETURNING *;

-- name: DeleteUserLoginCode :exec
DELETE
FROM user_login_codes
WHERE id = $1;

-- name: GetUserLoginCode :one
SELECT *
from user_login_codes
WHERE email = $1;

-- name: GetUserConfirmationToken :one
SELECT *
from user_confirmation_tokens
WHERE token = $1;

-- name: CreateUserConfirmationToken :one
INSERT INTO user_confirmation_tokens (email, token, expires_at)
VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteUserConfirmationToken :exec
DELETE
FROM user_confirmation_tokens
WHERE id = $1;