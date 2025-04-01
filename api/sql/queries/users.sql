-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY email;

-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2) RETURNING *;

-- name: UpdateUser :one
UPDATE users
set email    = $2,
    password = $3
WHERE id = $1 RETURNING *;

-- name: ConfirmUser :one
UPDATE users
set is_confirmed= true
WHERE email = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;