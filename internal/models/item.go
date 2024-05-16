package models

import "time"

type Item struct {
	ID        string    `json:"id" faker:"uuid_hyphenated"`
	Name      string    `json:"name" faker:"name"`
	Price     int       `json:"price" faker:"oneof: 15, 27, 61"`
	Quantity  int       `json:"quantity" faker:"oneof: 15, 27, 61"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
