package models

type User struct {
	ID       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email" binding:"required"`
	Password string `json:"password" db:"password" binding:"required"`
}

type Room struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name" binding:"required"`
	AdminID int    `json:"admin_id" db:"admin_id"`
}
