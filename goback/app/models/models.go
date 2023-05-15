package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type WSMessage struct {
	Email   string `json:"email"`
	Content string `json:"content"`
}

type ChatMessage struct {
	ID        int       `json:"id" db:"id"`
	RoomID    string    `json:"room_id" db:"room_id"`
	Content   string    `json:"content" db:"content"`
	Author    string    `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Room struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Author    string    `json:"author" db:"author"`
	IsOpen    bool      `json:"is_open" db:"is_open"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Quiz struct {
	ID        int       `json:"id" db:"id"`
	RoomID    string    `json:"room_id" db:"room_id"`
	Name      string    `json:"name" db:"name"`
	Author    string    `json:"author" db:"author"`
	Content   string    `json:"content" db:"content"`
	Variants  Variants  `json:"variants" db:"variants"`
	Answers   Answers   `json:"answers" db:"answers"`
	IsOpen    bool      `json:"is_open" db:"is_open"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Variant struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

type Variants []Variant

type Answer struct {
	VariantID int    `json:"variant_id"`
	Author    string `json:"author"`
}

type Answers []Answer
