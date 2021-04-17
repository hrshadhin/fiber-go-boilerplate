package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Book struct to describe book object.
type Book struct {
	ID        uuid.UUID  `db:"id" json:"id" validate:"uuid"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted bool       `db:"is_deleted" json:"is_deleted"`
	UserID    int        `db:"user_id" json:"user_id" validate:"required"`
	Title     string     `db:"title" json:"title" validate:"required,lte=255"`
	Author    string     `db:"author" json:"author" validate:"required,lte=255"`
	Status    int        `db:"status" json:"status" validate:"required,len=1"`
	Meta      Meta       `db:"meta" json:"meta" validate:"required,dive"`
}

// Meta struct to describe book attributes.
type Meta struct {
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Rating      int    `json:"rating" validate:"min=1,max=10"`
}

// Value make the Book Meta struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (b Meta) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan make the Book Meta struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (b *Meta) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &b)
}

func NewBook() *Book {
	return &Book{}
}
