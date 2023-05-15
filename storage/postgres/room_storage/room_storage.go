package room_storage

import (
	"errors"
	"fmt"
	"github.com/achekslime/core/models"
	"github.com/achekslime/core/storage/postgres"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type RoomStorage struct {
	db *sqlx.DB
}

func NewRoomStorage(db *sqlx.DB) *RoomStorage {
	return &RoomStorage{
		db: db,
	}
}

func (storage *RoomStorage) SaveRoom(room *models.Room) (int, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	query := fmt.Sprintf("INSERT INTO %s (name, admin_id, is_private) values ($1, $2) RETURNING id", postgres.RoomTableName)
	row := tx.QueryRow(query, room.Name, room.AdminID, room.IsPrivate)

	var id int
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (storage *RoomStorage) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s", postgres.RoomTableName)
	if err := storage.db.Get(&rooms, query); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (storage *RoomStorage) GetRoomsByAdminID(adminID int) ([]models.Room, error) {
	var rooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE admin_id=$1", postgres.RoomTableName)
	if err := storage.db.Get(&rooms, query, adminID); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (storage *RoomStorage) GetRoomByName(name string) ([]models.Room, error) {
	var rooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE name=$1", postgres.RoomTableName)
	if err := storage.db.Get(&rooms, query, name); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (storage *RoomStorage) GetAvailableRooms(userID int) ([]models.Room, error) {
	// public rooms.
	var publicRooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE is_private=$1", postgres.RoomTableName)
	if err := storage.db.Select(&publicRooms, query, false); err != nil {
		return nil, err
	}

	// private rooms.
	privateRooms, err := storage.getPrivateRooms(userID)
	if err != nil {
		return nil, err
	}

	if len(publicRooms) == 0 && len(privateRooms) == 0 {
		return nil, errors.New("sql: no rows in result set")
	}

	return append(publicRooms, privateRooms...), nil
}

func (storage *RoomStorage) getPrivateRooms(userID int) ([]models.Room, error) {
	// get ids from many-to-many table.
	var roomIDs []int
	getRoomsIdQuery := fmt.Sprintf("SELECT room_id FROM %s WHERE user_id=$1", postgres.AvailableRoomsTableName)
	if err := storage.db.Select(&roomIDs, getRoomsIdQuery, userID); err != nil {
		return nil, err
	}

	// if not exists private rooms.
	if len(roomIDs) == 0 {
		return nil, nil
	}

	// int to string.
	var ids []string
	for i := range roomIDs {
		number := roomIDs[i]
		text := strconv.Itoa(number)
		ids = append(ids, text)
	}

	// get private rooms list.
	var privateRooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", postgres.RoomTableName, strings.Join(ids, ","))
	if err := storage.db.Select(&privateRooms, query); err != nil {
		return nil, err
	}

	return privateRooms, nil
}

func (storage *RoomStorage) AddUsersToRoom(roomID int, userID int) error {
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

	query := fmt.Sprintf("INSERT INTO %s (room_id, user_id) "+
		"values ($1, $2)", postgres.AvailableRoomsTableName)

	if _, err = tx.Exec(query, roomID, userID); err != nil {
		return err
	}

	return err
}
