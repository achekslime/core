package postgres

import (
	"fmt"
	"github.com/achekslime/core/models"
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

	// type variation.
	var tableName string
	if room.Type == models.PrivateRoom {
		tableName = PrivateRoomsTableName
	} else {
		tableName = PublicRoomsTableName
	}

	query := fmt.Sprintf("INSERT INTO %s (name, admin_id) "+
		"values ($1, $2)", tableName)

	if _, err = tx.Exec(query, room.Name, room.AdminID); err != nil {
		return err
	}

	return err
}

func (storage *RoomStorage) GetAllRooms() ([]models.Room, error) {
	// private rooms.
	privateRooms, err := storage.getAllRooms(PrivateRoomsTableName)
	if err != nil {
		return nil, err
	}

	// public rooms.
	publicRooms, err := storage.getAllRooms(PublicRoomsTableName)
	if err != nil {
		return nil, err
	}

	return append(privateRooms, publicRooms...), nil
}

func (storage *RoomStorage) getAllRooms(tableName string) ([]models.Room, error) {
	var rooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	if err := storage.db.Select(&rooms, query); err != nil {
		return nil, err
	}

	// set room type.
	if tableName == PrivateRoomsTableName {
		for i := range rooms {
			rooms[i].Type = models.PrivateRoom
		}
	} else {
		for i := range rooms {
			rooms[i].Type = models.PublicRoom
		}
	}
	return rooms, nil
}

func (storage *RoomStorage) GetRoomsByAdminID(adminID int) ([]models.Room, error) {
	// private rooms.
	privateRooms, err := storage.getRoomsByAdminID(adminID, PrivateRoomsTableName)
	if err != nil {
		return nil, err
	}

	publicRooms, err := storage.getRoomsByAdminID(adminID, PublicRoomsTableName)
	if err != nil {
		return nil, err
	}

	return append(privateRooms, publicRooms...), nil
}

func (storage *RoomStorage) getRoomsByAdminID(adminID int, tableName string) ([]models.Room, error) {
	var rooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE admin_id=$1", tableName)
	if err := storage.db.Select(&rooms, query, adminID); err != nil {
		return nil, err
	}

	// set room type.
	if tableName == PrivateRoomsTableName {
		for i := range rooms {
			rooms[i].Type = models.PrivateRoom
		}
	} else {
		for i := range rooms {
			rooms[i].Type = models.PublicRoom
		}
	}
	return rooms, nil
}

func (storage *RoomStorage) GetAvailableRooms(userID int) ([]models.Room, error) {
	// get ids from many-to-many table.
	var roomIDs []int
	getRoomsIdQuery := fmt.Sprintf("SELECT room_id FROM %s WHERE user_id=$1", AvailableRoomsTableName)
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

	// get private rooms list.
	var privateRooms []models.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", PrivateRoomsTableName, strings.Join(ids, ","))
	if err := storage.db.Select(&privateRooms, query); err != nil {
		return nil, err
	}
	// set room type.
	for i := range privateRooms {
		privateRooms[i].Type = models.PrivateRoom
	}

	// public rooms.
	publicRooms, err := storage.getAllRooms(PublicRoomsTableName)
	if err != nil {
		return nil, err
	}

	return append(privateRooms, publicRooms...), nil
}
