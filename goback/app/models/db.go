package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

var UserColumns = []string{"id", "username", "created_at", "updated_at"}
var ChatColumns = []string{"id", "room_id", "content", "author", "created_at"}
var RoomColumns = []string{"id", "name", "author", "is_open", "created_at", "updated_at"}
var ActiveRoomUsersColumns = []string{"id", "room_id", "email", "connected", "active", "created_at"}
var QuizColumns = []string{"id", "room_id", "author", "name", "content", "variants", "answers", "is_open", "created_at", "updated_at"}

type RdMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Marshaller interface {
	Marshall() []byte
}

func MakeRdMessage(msgType string, data Marshaller) []byte {
	bts, _ := json.Marshal(RdMessage{
		Type: msgType,
		Data: string(data.Marshall()),
	})

	return bts
}

func MakeRdMessageStr(msgType, str string) []byte {
	bts, _ := json.Marshal(RdMessage{
		Type: msgType,
		Data: str,
	})

	return bts
}

func (cm *ChatMessage) Marshall() []byte {
	bts, _ := json.Marshal(cm)
	return bts
}

func (cm *Quiz) Marshall() []byte {
	bts, _ := json.Marshal(cm)
	return bts
}

// Value is used for storing into jsonb column
func (r Variants) Value() (driver.Value, error) {
	v := r
	if v == nil {
		v = make(Variants, 0)
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

// Scan is used for reading from jsonb column
func (r *Variants) Scan(v interface{}) error {
	bytes, ok := v.([]byte)
	if !ok {
		return errors.Errorf("Variants.Scan(): failed convert '%s' to []byte", v)
	}

	return json.Unmarshal(bytes, &r)
}

// Value is used for storing into jsonb column
func (r Answers) Value() (driver.Value, error) {
	v := r
	if v == nil {
		v = make(Answers, 0)
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

// Scan is used for reading from jsonb column
func (r *Answers) Scan(v interface{}) error {
	bytes, ok := v.([]byte)
	if !ok {
		return errors.Errorf("Answers.Scan(): failed convert '%s' to []byte", v)
	}

	return json.Unmarshal(bytes, &r)
}
