package storage

import (
	"github.com/achekslime/core/storage/models"
	"github.com/achekslime/core/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type UserStorage interface {
	SaveUser(*models.User) error
	GetAll() ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type RoomStorage interface {
	SaveRoom(room *models.Room) error
	GetAllRooms() ([]models.Room, error)
	GetRoomsByAdminID(userID int) ([]models.Room, error)
}

type Storage struct {
	userStorage UserStorage
	roomStorage RoomStorage
}

func NewStorage(db *sqlx.DB) (*Storage, error) {
	db, err := postgres.GetPostgresConnection()
	if err != nil {
		return nil, err
	}

	return &Storage{
		userStorage: postgres.NewUserStorage(db),
		roomStorage: postgres.NewRoomStorage(db),
	}, nil
}
