package model

import (
	"time"
)

// User struct to describe User object.
type User struct {
	ID        int        `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	IsActive  bool       `db:"is_active"`
	IsDeleted bool       `db:"is_deleted"`
	IsAdmin   bool       `db:"is_admin"`
	UserName  string     `db:"username"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	FirstName string     `db:"first_name"`
	LastName  string     `db:"last_name"`
}

func NewUser() *User {
	return &User{}
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUser struct {
	IsAdmin   bool   `json:"is_admin"`
	IsActive  bool   `json:"is_active"`
	UserName  string `json:"username" validate:"required,lte=50,gte=5"`
	Email     string `json:"email" validate:"required,email,lte=150"`
	Password  string `json:"password" validate:"required,lte=100,gte=10"`
	FirstName string `json:"first_name" validate:"required,lte=100"`
	LastName  string `json:"last_name" validate:"required,lte=100"`
}

type UpdateUser struct {
	IsAdmin   bool   `json:"is_admin"`
	IsActive  bool   `json:"is_active"`
	FirstName string `json:"first_name" validate:"required,lte=100"`
	LastName  string `json:"last_name" validate:"required,lte=100"`
}
