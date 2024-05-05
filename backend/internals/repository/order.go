package repository

import (
	"application/internals/dto"
	"application/internals/model"
	"application/pkg/pagination"
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IOrderRepository interface {
	FindAll(ctx context.Context, queries dto.FindAllRequest) ([]model.Order, uint, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	Create(ctx context.Context, data dto.CreateOrderRequest) error
	UpdateStatusByID(ctx context.Context, data dto.UpdateStatusOrderRequest) error
}

type order struct {
	db *sqlx.DB
}

func NewOrder(db *sqlx.DB) IOrderRepository {
	return &order{db: db}
}

func (r order) FindAll(ctx context.Context, queries dto.FindAllRequest) ([]model.Order, uint, error) {
	limit, offset := pagination.ToLimitOffset(*queries.Page, *queries.PerPage)
	sql := goqu.From("orders").
		Limit(limit).
		Offset(offset).
		Order(goqu.I("created_at").Desc())

	if queries.CreatedBy != nil {
		sql = sql.Where(goqu.Ex{
			"created_by": *queries.CreatedBy,
		})
	}

	if queries.Q != nil {
		sql = sql.Where(goqu.ExOr{
			"name":        goqu.Op{"like": *queries.Q},
			"description": goqu.Op{"like": *queries.Q},
		})
	}

	// Get
	query, _, err := sql.
		ToSQL()
	if err != nil {
		return nil, 0, err
	}

	var orders []model.Order
	rows, err := r.db.Queryx(query)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		var o model.Order
		err = rows.StructScan(&o)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	// ? Count
	query, _, err = sql.Select(goqu.COUNT("*")).ToSQL()
	if err != nil {
		return nil, 0, err
	}

	var total uint
	if err := r.db.Get(&total, query); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r order) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	query, _, err := goqu.
		From("orders").
		Where(goqu.Ex{
			"id": id,
		}).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var o model.Order
	err = r.db.QueryRowx(query).StructScan(&o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r order) Create(ctx context.Context, data dto.CreateOrderRequest) error {
	query, _, err := goqu.
		Insert("orders").
		Rows(goqu.Record{
			"id":          uuid.New(),
			"name":        data.Name,
			"description": data.Description,
			"status":      model.OrderDrafted,
			"created_by":  data.CreatedBy,
			"updated_by":  data.CreatedBy,
		}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (r order) UpdateStatusByID(ctx context.Context, data dto.UpdateStatusOrderRequest) error {
	now := time.Now()
	query, _, err := goqu.
		Update("orders").
		Set(goqu.Record{
			"updated_at": now,
			"updated_by": data.UpdatedBy,
			"note":       data.Note,
			"status":     data.Status,
		}).
		Where(goqu.Ex{"id": data.ID}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
