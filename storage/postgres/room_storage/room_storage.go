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

	query := fmt.Sprintf("INSERT INTO %s (name, admin_id, is_private) values ($1, $2, $3) RETURNING id", postgres.RoomTableName)
	row := tx.QueryRow(query, room.Name, room.AdminID, room.IsPrivate)

	var id int
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// add to available rooms.
	err = storage.AddUserToRoom(id, room.AdminID)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (storage *RoomStorage) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s", postgres.RoomTableName)
	if err := storage.db.Select(&rooms, query); err != nil {
		return nil, err
	}
	if len(rooms) == 0 {
		return nil, fmt.Errorf(postgres.ErrSqlNoRows)
	}
	return rooms, nil
}

func (storage *RoomStorage) GetRoomsByAdminID(adminID int) ([]models.Room, error) {
	var rooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE admin_id=$1", postgres.RoomTableName)
	if err := storage.db.Select(&rooms, query, adminID); err != nil {
		return nil, err
	}

	if len(rooms) == 0 {
		return nil, fmt.Errorf(postgres.ErrSqlNoRows)
	}
	return rooms, nil
}

func (storage *RoomStorage) GetRoomByName(name string) (*models.Room, error) {
	var room models.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE name=$1", postgres.RoomTableName)
	if err := storage.db.Get(&room, query, name); err != nil {
		return nil, err
	}
	return &room, nil
}

func (storage *RoomStorage) IsRoomAvailable(roomID int, userID int) (bool, error) {
	var roomIDs []int
	query := fmt.Sprintf("SELECT room_id FROM %s WHERE user_id=$1 && room_id=$2", postgres.AvailableRoomsTableName)
	if err := storage.db.Select(&roomIDs, query, userID, roomID); err != nil {
		return false, err
	}

	if len(roomIDs) == 0 {
		return false, fmt.Errorf(postgres.ErrSqlNoRows)
	}
	return true, nil
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
	getRoomsIdQuery := fmt.Sprintf("SELECT room_id FROM %s WHERE client_id=$1", postgres.AvailableRoomsTableName)
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

func (storage *RoomStorage) AddUserToRoom(roomID int, userID int) error {
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

	query := fmt.Sprintf("INSERT INTO %s (room_id, client_id) "+
		"values ($1, $2)", postgres.AvailableRoomsTableName)

	if _, err = tx.Exec(query, roomID, userID); err != nil {
		return err
	}

	return err
}

func (storage *RoomStorage) DelRoom(room models.Room) error {
	clientIDs, err := storage.getUsersID(room.ID)
	if err != nil {
		return err
	}

	// del from available rooms.
	for i := range clientIDs {
		err = storage.delAvailableRooms(room.ID, clientIDs[i])
		if err != nil {
			return err
		}
	}

	err = storage.delRoom(room.ID)
	if err != nil {
		return err
	}

	return nil
}

func (storage *RoomStorage) getUsersID(roomID int) ([]int, error) {
	var clientIDs []int
	getClientsIDs := fmt.Sprintf("SELECT client_id FROM %s WHERE room_id=%d", postgres.AvailableRoomsTableName, roomID)
	if err := storage.db.Select(&clientIDs, getClientsIDs); err != nil {
		return nil, err
	}
	return clientIDs, nil
}

func (storage *RoomStorage) delAvailableRooms(roomID int, userID int) error {
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

	query := fmt.Sprintf("DELETE FROM %s WHERE room_id=%d AND client_id=%d", postgres.AvailableRoomsTableName, roomID, userID)

	if _, err = tx.Exec(query); err != nil {
		return err
	}

	return err
}

func (storage *RoomStorage) delRoom(roomID int) error {
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

	query := fmt.Sprintf("DELETE FROM %s WHERE id=%d", postgres.RoomTableName, roomID)

	if _, err = tx.Exec(query); err != nil {
		return err
	}

	return err
}

func (storage *RoomStorage) ClearAll() error {
	err := storage.clearAvailableRooms()
	if err != nil {
		return err
	}

	err = storage.clearRooms()
	if err != nil {
		return err
	}

	err = storage.clearUsers()
	if err != nil {
		return err
	}

	return nil
}

func (storage *RoomStorage) clearAvailableRooms() error {
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
	query := fmt.Sprintf("DELETE FROM %s", postgres.AvailableRoomsTableName)
	if _, err = tx.Exec(query); err != nil {
		return err
	}
	return err
}

func (storage *RoomStorage) clearRooms() error {
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
	query := fmt.Sprintf("DELETE FROM %s", postgres.RoomTableName)
	if _, err = tx.Exec(query); err != nil {
		return err
	}
	return err
}

func (storage *RoomStorage) clearUsers() error {
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
	query := fmt.Sprintf("DELETE FROM %s", postgres.UserTableName)
	if _, err = tx.Exec(query); err != nil {
		return err
	}
	return err
}
