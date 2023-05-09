package replatform

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/asiman161/re-platform/app/models"
	"github.com/asiman161/re-platform/app/room"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}} // use default options

func (i *Implementation) CreateRoom(w http.ResponseWriter, r *http.Request) {
	req := models.Room{}

	u := r.Header.Get("Email")
	fmt.Println(r.Header, u)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.PlainText(w, r, "can't decode room")
	}

	newRoom, err := i.store.CreateRoom(r.Context(), req.Name, extractAuthor(r))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.PlainText(w, r, "can't create room")
	}

	render.JSON(w, r, newRoom)
}

func (i *Implementation) CloseRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "ID")

	err := i.store.CloseRoom(r.Context(), roomID, extractAuthor(r))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			render.PlainText(w, r, fmt.Sprintf("can't find room: %s", roomID))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.PlainText(w, r, "can't close room")
	}

	render.PlainText(w, r, "OK")
}

func (i *Implementation) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := i.store.GetOpenRooms(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.PlainText(w, r, "can't get rooms")
		return
	}

	render.JSON(w, r, rooms)
}

func (i *Implementation) Room(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "ID")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.PlainText(w, r, "can't open socket connection")
		log.Print("upgrade:", err)
		return
	}

	newRoom := room.New(c, i.store, roomID)
	newRoom.Connect()
}

func (i *Implementation) GetMessages(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "ID")

	messages, err := i.store.GetMessages(r.Context(), roomID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.PlainText(w, r, "can't save message")
	}

	render.JSON(w, r, messages)
}
