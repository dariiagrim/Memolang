package models

import (
	"time"
)

type User struct {
	Id          int64     `gorm:"primaryKey" json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Password    string    `json:"password"`
	Username    string    `json:"username" gorm:"unique"`
	DateOfBirth time.Time `json:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at"`
}
