package room

import (
	"log"

	"github.com/gorilla/websocket"
)

type Room struct {
	conn *websocket.Conn
}

func New(conn *websocket.Conn) *Room {
	return &Room{conn: conn}
}

func (r *Room) Connect() {
	defer r.conn.Close()
	for {
		mt, message, err := r.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		err = r.conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
