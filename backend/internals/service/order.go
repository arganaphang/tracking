package service

import (
	"application/internals/dto"
	"application/internals/model"
	"application/internals/repository"
	"context"

	"github.com/google/uuid"
)

type IOrderService interface {
	FindAll(ctx context.Context, queries dto.FindAllRequest) ([]model.Order, uint, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	Create(ctx context.Context, data dto.CreateOrderRequest) error
	UpdateStatusByID(ctx context.Context, data dto.UpdateStatusOrderRequest) error
}

type order struct {
	repositories repository.Repositories
}

func NewOrder(repositories repository.Repositories) IOrderService {
	return &order{repositories: repositories}
}

func (s order) FindAll(ctx context.Context, queries dto.FindAllRequest) ([]model.Order, uint, error) {
	return s.repositories.Order.FindAll(ctx, queries)
}

func (s order) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	return s.repositories.Order.FindByID(ctx, id)
}

func (s order) Create(ctx context.Context, data dto.CreateOrderRequest) error {
	return s.repositories.Order.Create(ctx, data)
}

func (s order) UpdateStatusByID(ctx context.Context, data dto.UpdateStatusOrderRequest) error {
	return s.repositories.Order.UpdateStatusByID(ctx, data)
}
