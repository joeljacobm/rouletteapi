package controllers

import (
	"encoding/json"
	"net/http"
	"rouletteapi/configs"
	"rouletteapi/models"
	"time"

	"github.com/gorilla/mux"
)

type Room struct {
	Room models.RoomService
	r    *mux.Router
}

func NewRoom(room models.RoomService, r *mux.Router) *Room {
	return &Room{
		Room: room,
		r:    r,
	}
}

func (ro *Room) RoomHandler(w http.ResponseWriter, r *http.Request) {

	var room models.Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}
	roomID := GenerateHash()

	room = configs.GetRouletteVariantMap(room.VariantType)

	room.ID = roomID
	room.Created = time.Now()
	room.CurrentRound = 0

	err = ro.Room.Create(room)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}

	resp := struct {
		message string
		Data    models.Room
	}{
		message: "Successfully created the room",
		Data:    room,
	}
	writeJSON(w, resp)
}

