package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	OrderDrafted    Status = "DRAFTED"
	OrderProcessing Status = "PROCESSING"
	OrderRejected   Status = "REJECTED"
	OrderApproved   Status = "APPROVED"
)

var (
	orders = [...]Status{OrderDrafted, OrderProcessing, OrderRejected, OrderApproved}
)

func OrderStatusContains(status string) bool {
	for _, s := range orders {
		if strings.EqualFold(string(s), status) {
			return true
		}
	}

	return false
}

type Order struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Note        *string   `json:"note" db:"note"`
	Status      Status    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy   uuid.UUID `json:"updated_by" db:"updated_by"`
}
