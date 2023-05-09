package models

import (
	"net/http"
	"time"
)

var UserColumns = []string{"id", "username", "created_at", "updated_at"}
var ChatColumns = []string{"id", "chat_id", "content", "author", "created_at"}
var RoomColumns = []string{"id", "name", "author", "is_open", "created_at", "updated_at"}

type User struct {
	Id        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type WSMessage struct {
	Email   string `json:"email"`
	Content string `json:"content"`
}

type ChatMessage struct {
	Id        int       `json:"id" db:"id"`
	ChatID    string    `json:"chat_id" db:"chat_id"`
	Content   string    `json:"content" db:"content"`
	Author    string    `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (cm *ChatMessage) Bind(r *http.Request) error {
	return nil
}

type Room struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Author    string    `json:"author" db:"author"`
	IsOpen    bool      `json:"is_open" db:"is_open"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
