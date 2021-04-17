package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
	"time"
)

type BookRepo struct {
	db *database.DB
}

func NewBookRepo(db *database.DB) BookRepository {
	return &BookRepo{db}
}

func (repo *BookRepo) Create(b *model.Book) error {
	query := `INSERT INTO book VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := repo.db.Exec(query, b.ID, time.Now(), time.Now(), false, b.UserID, b.Title, b.Author, b.Status, b.Meta)
	return err
}

func (repo *BookRepo) All(limit int, offset uint) ([]*model.Book, error) {
	var books []*model.Book
	query := `SELECT * FROM book WHERE is_deleted = FALSE`
	var err error

	if limit > 0 && offset >= 0 {
		query = fmt.Sprintf("%s LIMIT $1 OFFSET $2", query)
		err = repo.db.Select(&books, query, limit, offset)
	} else {
		err = repo.db.Select(&books, query)
	}

	return books, err
}

func (repo *BookRepo) Get(ID uuid.UUID) (*model.Book, error) {
	book := model.Book{}
	query := `SELECT * FROM book WHERE id = $1 AND is_deleted = FALSE`
	err := repo.db.Get(&book, query, ID)

	return &book, err
}

func (repo *BookRepo) Update(ID uuid.UUID, b *model.Book) error {
	query := `UPDATE book SET updated_at = $2, title = $3, author = $4, status = $5, meta = $6 WHERE id = $1`
	_, err := repo.db.Exec(query, ID, time.Now(), b.Title, b.Author, b.Status, b.Meta)
	return err
}

func (repo *BookRepo) Delete(ID uuid.UUID) error {
	query := `UPDATE book SET is_deleted = TRUE, updated_at = $2 WHERE id = $1`
	_, err := repo.db.Exec(query, ID, time.Now())
	return err
}
