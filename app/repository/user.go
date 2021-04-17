package repository

import (
	"database/sql"
	"fmt"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
	"time"
)

type UserRepo struct {
	db *database.DB
}

func NewUserRepo(db *database.DB) UserRepository {
	return &UserRepo{db}
}

func (repo *UserRepo) Create(u *model.CreateUser) error {
	query := `INSERT INTO "user" (is_active, is_admin, username, email, password, first_name, last_name) VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err := repo.db.Exec(query, u.IsActive, u.IsAdmin, u.UserName, u.Email, u.Password, u.FirstName, u.LastName)
	return err
}

func (repo *UserRepo) All(limit int, offset uint) ([]*model.User, error) {
	var users []*model.User
	query := `SELECT * FROM "user" WHERE is_deleted=FALSE`
	var err error

	if limit > 0 && offset >= 0 {
		query = fmt.Sprintf("%s LIMIT $1 OFFSET $2", query)
		err = repo.db.Select(&users, query, limit, offset)
	} else {
		err = repo.db.Select(&users, query)
	}

	return users, err
}

func (repo *UserRepo) GetByUsername(username string) (*model.User, error) {
	user := model.User{}
	query := `SELECT * FROM "user" WHERE username = $1 AND is_deleted = FALSE`
	err := repo.db.Get(&user, query, username)

	return &user, err
}

func (repo *UserRepo) Get(ID int) (*model.User, error) {
	user := model.User{}
	query := `SELECT * FROM "user" WHERE id = $1 AND is_deleted = FALSE`
	err := repo.db.Get(&user, query, ID)

	return &user, err
}

func (repo *UserRepo) Exists(username, email string) (bool, error) {
	findByUser := len(username)
	findByEmail := len(email)
	if findByUser <= 0 || findByEmail <= 0 {
		return false, nil
	}
	var placeHolderValues []interface{}
	query := `SELECT 1 FROM "user"`
	if findByUser > 0 && findByEmail > 0 {
		query = fmt.Sprintf("%s WHERE username = $1 OR email = $2", query)
		placeHolderValues = append(placeHolderValues, username, email)
	} else {
		if findByUser > 0 {
			query = fmt.Sprintf("%s WHERE username = $1", query)
			placeHolderValues = append(placeHolderValues, username)
		} else {
			query = fmt.Sprintf("%s WHERE email = $1", query)
			placeHolderValues = append(placeHolderValues, email)
		}
	}

	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := repo.db.QueryRow(query, placeHolderValues...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func (repo *UserRepo) Update(ID int, u *model.UpdateUser) error {
	query := `UPDATE "user" SET updated_at = $2, is_active = $3, is_admin = $4, first_name = $5, last_name = $6 WHERE id = $1`
	_, err := repo.db.Exec(query, ID, time.Now(), u.IsActive, u.IsAdmin, u.FirstName, u.LastName)
	return err
}

func (repo *UserRepo) Delete(ID int) error {
	query := `UPDATE "user" SET is_deleted = TRUE,updated_at = $2 WHERE id = $1`
	_, err := repo.db.Exec(query, ID, time.Now())
	return err
}
