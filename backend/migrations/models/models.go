package models

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type User struct {
	Uuid        string      `gorm:"primaryKey" json:"uuid"`
	FirstName   null.String `json:"first_name"`
	LastName    null.String `json:"last_name"`
	Email       string      `json:"email"`
	Username    string      `json:"username" gorm:"unique"`
	DateOfBirth null.Time   `json:"date_of_birth"`
	CreatedAt   time.Time   `json:"created_at"`
	Avatar      null.String `json:"avatar"`
	Points      int         `json:"points"`
}

type Topic struct {
	Id        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
	General   bool      `json:"general"`
}

type Word struct {
	Id              int64     `gorm:"primaryKey" json:"id"`
	EnglishSpelling string    `json:"english_spelling"`
	FrenchSpelling  string    `json:"french_spelling"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       null.Time `json:"deleted_at"`
}

type TopicsWord struct {
	Id        int64     `gorm:"primaryKey" json:"id"`
	TopicId   int64     `json:"topic_id"`
	WordId    int64     `json:"word_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

type TopicsUser struct {
	Id        int64     `gorm:"primaryKey" json:"id"`
	TopicId   int64     `json:"topic_id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}
