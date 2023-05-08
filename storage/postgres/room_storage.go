package postgres

import (
	"fmt"
	"github.com/achekslime/core/models"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

const AvailableRoomsTableName = "available_rooms"

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
	query := fmt.Sprintf("SELECT * FROM %s", storage.tableName)

	var rooms []models.Room
	if err := storage.db.Get(&rooms, query); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (storage *RoomStorage) GetRoomsByAdminID(adminID int) ([]models.Room, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE admin_id=$1 ", storage.tableName)

	var rooms []models.Room
	if err := storage.db.Select(&rooms, query, adminID); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (storage *RoomStorage) GetAvailableRooms(userID int) ([]models.Room, error) {
	// get ids from many-to-many table.
	getRoomsIdQuery := fmt.Sprintf("SELECT room_id FROM %s WHERE user_id=$1 ", AvailableRoomsTableName)
	var roomIDs []int
	if err := storage.db.Select(&roomIDs, getRoomsIdQuery, userID); err != nil {
		return nil, err
	}

	// int to string.
	var ids []string
	for i := range roomIDs {
		number := roomIDs[i]
		text := strconv.Itoa(number)
		ids = append(ids, text)
	}

	// grt rooms list.
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", storage.tableName, strings.Join(ids, ","))
	var rooms []models.Room
	if err := storage.db.Select(&rooms, query); err != nil {
		return nil, err
	}

	return rooms, nil
}
