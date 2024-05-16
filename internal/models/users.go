package models

import "time"

type User struct {
	ID        string    `json:"id" faker:"uuid_hyphenated"`
	Email     string    `json:"email" gorm:"unique" faker:"email"`
	Password  string    `json:"password" faker:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
