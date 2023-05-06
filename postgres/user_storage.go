package postgres

import (
	"fmt"
	"github.com/achekslime/core/models"
	"github.com/jmoiron/sqlx"
)

const (
	TableName = "client"
)

type UserStorage struct {
	db        *sqlx.DB
	tableName string
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db:        db,
		tableName: TableName,
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

	query := fmt.Sprintf("INSERT INTO %s (email, password) "+
		"values ($1, $2)", storage.tableName)

	if _, err = tx.Exec(query, user.Email, user.Password); err != nil {
		return err
	}

	return err
}

func (storage *UserStorage) GetUserByEmail(email string) (*models.User, error) {
	query := fmt.Sprintf("SELECT email, password FROM %s WHERE email=$1 ", storage.tableName)

	var user models.User
	if err := storage.db.Get(&user, query, email); err != nil {
		return nil, err
	}

	return &user, nil
}
