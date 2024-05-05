package dto

import (
	"application/internals/model"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"-"`
}

type CreateOrderResponse struct {
	Message string `json:"message"`
}

type UpdateStatusOrderRequest struct {
	ID        uuid.UUID `json:"-"`
	Note      *string   `json:"note"`
	Status    string    `json:"-"`
	UpdatedBy uuid.UUID `json:"-"`
}

type UpdateStatusOrderResponse struct {
	Message string `json:"message"`
}

type FindAllRequest struct {
	Page      *uint      `form:"page"`
	PerPage   *uint      `form:"per_page"`
	Q         *string    `form:"q"`
	CreatedBy *uuid.UUID `form:"-"`
}

type FindAllResponse struct {
	Message string        `json:"message"`
	Data    []model.Order `json:"data"`
	Meta    Meta          `json:"meta"`
}

type FindByIDResponse struct {
	Message string       `json:"message"`
	Data    *model.Order `json:"data"`
}
