package storage

import "fmt"

func redisChatID(id string) string {
	return fmt.Sprintf("chat:%s", id)
}
