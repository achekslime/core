package storage

import (
	"github.com/achekslime/core/models"
	"github.com/achekslime/core/storage/postgres"
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
		UserStorage: postgres.NewUserStorage(db),
		RoomStorage: postgres.NewRoomStorage(db),
	}, nil
}
