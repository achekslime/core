package models

import (
	"github.com/achekslime/core/models"
)

type User struct {
	ID       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Room struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	AdminID int    `json:"admin_id" db:"admin_id"`
}

func UserToDbUser(user *models.User) *User {
	return &User{
		Email:    user.Email,
		Password: user.Password,
	}
}

func RoomToDbRoom(room *models.Room, adminID int) *Room {
	return &Room{
		Name:    room.Name,
		AdminID: adminID,
	}
}
