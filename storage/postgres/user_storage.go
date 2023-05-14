package postgres

import (
	"fmt"
	"github.com/achekslime/core/models"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UserStorage struct {
	db        *sqlx.DB
	tableName string
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db:        db,
		tableName: UserTableName,
	}
}

func (storage *UserStorage) SaveUser(user *models.User) error {
	tx, err := storage.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	query := fmt.Sprintf("INSERT INTO %s (name, email, password) "+
		"values ($1, $2, $3)", storage.tableName)

	if _, err = tx.Exec(query, user.Name, user.Email, user.Password); err != nil {
		return err
	}

	return err
}

func (storage *UserStorage) GetAll() ([]models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", storage.tableName)

	var user []models.User
	if err := storage.db.Get(&user, query); err != nil {
		return nil, err
	}

	return user, nil
}

func (storage *UserStorage) GetUserByEmail(email string) (*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 ", storage.tableName)

	var user models.User
	if err := storage.db.Get(&user, query, email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (storage *UserStorage) GetUserByName(name string) (*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE name=$1 ", storage.tableName)

	var user models.User
	if err := storage.db.Get(&user, query, name); err != nil {
		return nil, err
	}

	return &user, nil
}
