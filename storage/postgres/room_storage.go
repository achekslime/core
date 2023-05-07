package postgres

import (
	"fmt"
	"github.com/achekslime/core/storage/models"
	"github.com/jmoiron/sqlx"
)

type RoomStorage struct {
	db        *sqlx.DB
	tableName string
}

func NewRoomStorage(db *sqlx.DB) *RoomStorage {
	return &RoomStorage{
		db:        db,
		tableName: RoomTableName,
	}
}

func (storage *RoomStorage) SaveRoom(room *models.Room) error {
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

	query := fmt.Sprintf("INSERT INTO %s (name, admin_id) "+
		"values ($1, $2)", storage.tableName)

	if _, err = tx.Exec(query, room.Name, room.AdminID); err != nil {
		return err
	}

	return err
}

func (storage *RoomStorage) GetAllRooms() ([]models.Room, error) {
	query := fmt.Sprintf("SELECT id, name, admin_id FROM %s", storage.tableName)

	var rooms []models.Room
	if err := storage.db.Get(&rooms, query); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (storage *RoomStorage) GetRoomsByAdminID(userID int) ([]models.Room, error) {
	query := fmt.Sprintf("SELECT id, name, admin_id FROM %s WHERE admin_id=$1 ", storage.tableName)

	var rooms []models.Room
	if err := storage.db.Get(&rooms, query, userID); err != nil {
		return nil, err
	}

	return rooms, nil
}
