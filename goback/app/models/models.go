package models

import "time"

var UserColumns = []string{"id", "username", "created_at", "updated_at"}

type User struct {
	Id        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
