package service

import (
	"application/internals/dto"
	"application/internals/model"
	"application/internals/repository"
	"context"

	"github.com/google/uuid"
)

type IUserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, data dto.RegisterRequest) error
	Active(ctx context.Context, id uuid.UUID) error
}

type user struct {
	repositories repository.Repositories
}

func NewUser(repositories repository.Repositories) IUserService {
	return &user{repositories: repositories}
}

func (s user) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repositories.User.GetByID(ctx, id)
}

func (s user) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repositories.User.GetByEmail(ctx, email)
}

func (s user) Create(ctx context.Context, data dto.RegisterRequest) error {
	return s.repositories.User.Create(ctx, data)
}

func (s user) Active(ctx context.Context, id uuid.UUID) error {
	return s.repositories.User.Active(ctx, id)
}
