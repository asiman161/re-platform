package room

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/asiman161/re-platform/app/models"
	"github.com/asiman161/re-platform/storage"
	"github.com/gorilla/websocket"
)

type Room struct {
	conn   *websocket.Conn
	store  storage.Storager
	roomID string
}

func New(conn *websocket.Conn, store storage.Storager, roomID string) *Room {
	return &Room{conn: conn, store: store, roomID: roomID}
}

func (r *Room) Connect() {
	defer r.conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := r.conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}

				wsMessage := models.WSMessage{}
				err = json.Unmarshal(message, &wsMessage)
				if err != nil {
					return
				}

				msg := models.ChatMessage{
					RoomID:    r.roomID,
					Content:   wsMessage.Content,
					Author:    wsMessage.Email,
					CreatedAt: time.Now(),
				}

				_, err = r.store.WriteChatMessage(context.Background(), msg)
				if err != nil {
					cancel()
					log.Println("close chat, err: ", err)
					return
				}
			}

		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ch, err := r.store.SubscribeMessages(ctx, r.roomID)
		if err != nil {
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-ch:
				if !ok {
					return
				}

				err = r.conn.WriteJSON(v)
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}
	}()

	wg.Wait()
}
