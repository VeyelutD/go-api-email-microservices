package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/VeyelutD/go-api-microservice/internal/db"
	"github.com/VeyelutD/go-api-microservice/internal/email"
	"github.com/VeyelutD/go-api-microservice/internal/otp"
	"github.com/VeyelutD/go-api-microservice/internal/tokens"
	"github.com/VeyelutD/go-api-microservice/internal/users"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Service struct {
	queries *db.Queries
	db      *pgx.Conn
	email   *email.Service
	users   *users.Service
}

func NewService(queries *db.Queries, db *pgx.Conn, emailService *email.Service, userService *users.Service) *Service {
	return &Service{
		queries: queries,
		db:      db,
		email:   emailService,
		users:   userService,
	}
}

func (s *Service) GetUserOTP(ctx context.Context, email string) (*db.UserLoginCode, error) {
	user, err := s.queries.GetUserLoginCode(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserOTPNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) DeleteUserOTP(ctx context.Context, id int64) error {
	err := s.queries.DeleteUserLoginCode(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateUserOTP(ctx context.Context, email string) (*db.UserLoginCode, error) {
	code, err := otp.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("error generating otp: %w", err)
	}
	userLoginCode, err := s.queries.CreateUserLoginCode(ctx, db.CreateUserLoginCodeParams{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return nil, err
	}
	return &userLoginCode, nil
}
func (s *Service) VerifyOTPAndGetUser(ctx context.Context, email, code string) (*db.User, error) {
	user, err := s.users.GetOneByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !user.IsConfirmed {
		return nil, errors.New("user is not confirmed")
	}
	userLoginCode, err := s.GetUserOTP(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if code != userLoginCode.Code {
		return nil, ErrWrongOTP
	}
	if err = s.DeleteUserOTP(ctx, userLoginCode.ID); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) CreateUserConfirmationToken(ctx context.Context, email string) (*db.UserConfirmationToken, error) {
	token, err := tokens.GenerateConfirmationToken(32)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}
	userLoginCode, err := s.queries.CreateUserConfirmationToken(ctx, db.CreateUserConfirmationTokenParams{
		Email:     email,
		Token:     token,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24), Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &userLoginCode, nil
}

func (s *Service) GetUserConfirmationToken(ctx context.Context, token string) (db.UserConfirmationToken, error) {
	user, err := s.queries.GetUserConfirmationToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.UserConfirmationToken{}, ErrUserConfirmationTokenNotFound
		}
		return db.UserConfirmationToken{}, err
	}
	return user, nil
}

func (s *Service) ConfirmUserAndDeleteConfirmationToken(ctx context.Context, email string, tokenID int64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)
	_, err = qtx.ConfirmUser(ctx, email)
	if err := qtx.DeleteUserConfirmationToken(ctx, tokenID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
