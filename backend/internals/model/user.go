package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin     Role = "ADMIN"
	RoleApplicant Role = "APPLICANT"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Role      Role      `json:"role" db:"role"`
	Password  string    `json:"-" db:"password"`
	IsActive  bool      `json:"-" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (u User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
