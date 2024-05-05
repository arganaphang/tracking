package repository

import (
	"application/internals/dto"
	"application/internals/model"
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type IUserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, data dto.RegisterRequest) error
	Active(ctx context.Context, id uuid.UUID) error
}

type user struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) IUserRepository {
	return &user{db: db}
}

func (r user) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query, _, err := goqu.
		From("users").
		Where(goqu.Ex{
			"id": id,
		}).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = r.db.QueryRowx(query).StructScan(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r user) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query, _, err := goqu.
		From("users").
		Where(goqu.Ex{
			"email": email,
		}).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = r.db.QueryRowx(query).StructScan(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r user) Create(ctx context.Context, data dto.RegisterRequest) error {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	password := string(hashedByte)
	query, _, err := goqu.
		Insert("users").
		Rows(goqu.Record{
			"id":       uuid.New(),
			"name":     data.Name,
			"email":    data.Email,
			"role":     model.RoleApplicant,
			"password": password,
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

func (r user) Active(ctx context.Context, id uuid.UUID) error {
	query, _, err := goqu.
		Update("users").
		Set(goqu.Record{
			"is_active": 1,
		}).
		Where(goqu.Ex{"id": id}).
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
