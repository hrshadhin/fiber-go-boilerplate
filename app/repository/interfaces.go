package repository

import (
	"github.com/google/uuid"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
)

type UserRepository interface {
	Create(b *model.CreateUser) error
	All(limit int, offset uint) ([]*model.User, error)
	Get(ID int) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Exists(username, email string) (bool, error)
	Update(ID int, user *model.UpdateUser) error
	Delete(ID int) error
}

type BookRepository interface {
	Create(b *model.Book) error
	All(limit int, offset uint) ([]*model.Book, error)
	Get(ID uuid.UUID) (*model.Book, error)
	Update(ID uuid.UUID, b *model.Book) error
	Delete(ID uuid.UUID) error
}
