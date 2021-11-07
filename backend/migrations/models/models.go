package models

import (
	"gopkg.in/guregu/null.v4"
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

type Topic struct {
	Id          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   null.Time `json:"deleted_at"`
	General     bool      `json:"general"`
}

type Word struct {
	Id              int64  `gorm:"primaryKey" json:"id"`
	EnglishSpelling string `json:"english_spelling"`
	FrenchSpelling  string `json:"french_spelling"`
	Definition      string `json:"definition"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   null.Time `json:"deleted_at"`
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
	UserId    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}
