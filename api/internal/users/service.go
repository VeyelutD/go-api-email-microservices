package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/VeyelutD/go-api-microservice/internal/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

func (s *Service) GetOneByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := s.queries.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) Create(ctx context.Context, email string) (*db.User, error) {
	newUser, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	return &newUser, nil
}
