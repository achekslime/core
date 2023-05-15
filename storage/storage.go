package storage

import (
	"github.com/achekslime/core/models"
	"github.com/achekslime/core/storage/postgres"
	"github.com/achekslime/core/storage/postgres/room_storage"
	"github.com/achekslime/core/storage/postgres/user_storage"
)

type UserStorage interface {
	SaveUser(*models.User) error
	GetAll() ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
}

type RoomStorage interface {
	SaveRoom(room *models.Room) (int, error)
	GetAllRooms() ([]models.Room, error)
	GetRoomsByAdminID(userID int) ([]models.Room, error)
	GetRoomByName(name string) ([]models.Room, error)
	GetAvailableRooms(userID int) ([]models.Room, error)
	AddUsersToRoom(roomID int, userID int) error
}

type Storage struct {
	UserStorage UserStorage
	RoomStorage RoomStorage
}

func NewStorage() (*Storage, error) {
	db, err := postgres.GetPostgresConnection()
	if err != nil {
		return nil, err
	}

	return &Storage{
		UserStorage: user_storage.NewUserStorage(db),
		RoomStorage: room_storage.NewRoomStorage(db),
	}, nil
}
