package models

type User struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Room struct {
	Name  string `json:"name" binding:"required"`
	Admin *User  `json:"admin" binding:"required"`
}
