package replatform

import (
	"fmt"
	"log"
	"net/http"

	"github.com/asiman161/re-platform/app/room"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}} // use default options

func (i *Implementation) Room(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "ID")

	fmt.Println("test log roomID:", roomID)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.PlainText(w, r, "can't open socket connection")
		log.Print("upgrade:", err)
		return
	}

	newRoom := room.New(c)
	newRoom.Connect()
}
