// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID          int64            `json:"id"`
	Email       string           `json:"email"`
	Password    pgtype.Text      `json:"password"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	IsConfirmed bool             `json:"is_confirmed"`
}

type UserConfirmationToken struct {
	ID        int64            `json:"id"`
	Email     string           `json:"email"`
	Token     string           `json:"token"`
	ExpiresAt pgtype.Timestamp `json:"expires_at"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type UserLoginCode struct {
	ID        int64            `json:"id"`
	Email     string           `json:"email"`
	Code      string           `json:"code"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}
