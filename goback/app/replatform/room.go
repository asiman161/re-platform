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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "can't decode room")
		return
	}

	newRoom, err := i.store.CreateRoom(r.Context(), req.Name, extractAuthor(r))
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "can't create room")
		return
	}

	render.JSON(w, r, newRoom)
}

func (i *Implementation) CloseRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_ID")

	err := i.store.CloseRoom(r.Context(), roomID, extractAuthor(r))
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, r, http.StatusNotFound, fmt.Sprintf("can't find room: %s", roomID))
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
		render.PlainText(w, r, "can't get rooms")
		writeError(w, r, http.StatusInternalServerError, "can't get rooms")
		return
	}

	render.JSON(w, r, rooms)
}

func (i *Implementation) GetRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_ID")

	foundRoom, err := i.store.GetRoom(r.Context(), roomID)
	if err != nil {
		render.PlainText(w, r, "can't get rooms")
		writeError(w, r, http.StatusInternalServerError, "can't get rooms")
		return
	}

	render.JSON(w, r, foundRoom)
}

func (i *Implementation) Room(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_ID")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		writeError(w, r, http.StatusInternalServerError, "can't open socket connection")
		return
	}

	newRoom := room.New(c, i.store, roomID)
	newRoom.Connect()
}

func (i *Implementation) GetMessages(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_ID")

	messages, err := i.store.GetMessages(r.Context(), roomID)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "can't save message")
		return
	}

	render.JSON(w, r, messages)
}
