package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

var UserColumns = []string{"id", "username", "created_at", "updated_at"}
var ChatColumns = []string{"id", "room_id", "content", "author", "created_at"}
var RoomColumns = []string{"id", "name", "author", "is_open", "created_at", "updated_at"}
var PoolColumns = []string{"id", "room_id", "author", "content", "variants", "answers", "is_open", "created_at", "updated_at"}

type RdMessage struct {
	Type string      `json:"type"`
	Data ChatMessage `json:"data"`
}

func MakeRdMessage(msgType string, data ChatMessage) []byte {
	bts, _ := json.Marshal(RdMessage{
		Type: msgType,
		Data: data,
	})

	return bts
}

// Value is used for storing region into jsonb column
func (r Variants) Value() (driver.Value, error) {
	// glorious sql.Balancer driver crutch
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

// Scan is used for reading Location from jsonb column
func (r *Variants) Scan(v interface{}) error {
	bytes, ok := v.([]byte)
	if !ok {
		return errors.Errorf("Variants.Scan(): failed convert '%s' to []byte", v)
	}

	return json.Unmarshal(bytes, &r)
}

// Value is used for storing region into jsonb column
func (r Answers) Value() (driver.Value, error) {
	// glorious sql.Balancer driver crutch
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

// Scan is used for reading Location from jsonb column
func (r *Answers) Scan(v interface{}) error {
	bytes, ok := v.([]byte)
	if !ok {
		return errors.Errorf("Answers.Scan(): failed convert '%s' to []byte", v)
	}

	return json.Unmarshal(bytes, &r)
}
