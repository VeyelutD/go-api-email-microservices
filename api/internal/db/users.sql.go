// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const confirmUser = `-- name: ConfirmUser :one
UPDATE users
set is_confirmed= true
WHERE email = $1 RETURNING id, email, password, created_at, is_confirmed
`

func (q *Queries) ConfirmUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, confirmUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.IsConfirmed,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2) RETURNING id, email, password, created_at, is_confirmed
`

type CreateUserParams struct {
	Email    string      `json:"email"`
	Password pgtype.Text `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.IsConfirmed,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, email, password, created_at, is_confirmed
FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.IsConfirmed,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, password, created_at, is_confirmed
FROM users
ORDER BY email
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.IsConfirmed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
set email    = $2,
    password = $3
WHERE id = $1 RETURNING id, email, password, created_at, is_confirmed
`

type UpdateUserParams struct {
	ID       int64       `json:"id"`
	Email    string      `json:"email"`
	Password pgtype.Text `json:"password"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser, arg.ID, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.IsConfirmed,
	)
	return i, err
}
