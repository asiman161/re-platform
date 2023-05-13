package storage

import (
	"fmt"
	"strings"
)

func redisRoomID(id string) string {
	return fmt.Sprintf("room:%s", id)
}

func suffixReturning(columns []string) string {
	return fmt.Sprintf("RETURNING %s", strings.Join(columns, ","))
}
