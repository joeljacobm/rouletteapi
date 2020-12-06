package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"rouletteapi/configs"
	"rouletteapi/models"
	"time"

	"github.com/gorilla/mux"
)

type Room struct {
	Room models.RoomService
}

func NewRoom(room models.RoomService) *Room {
	return &Room{
		Room: room,
	}
}

func (ro *Room) RoomHandler(w http.ResponseWriter, r *http.Request) {

	var room models.Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		writeErrorWithMsg(w, err)
		return
	}
	roomID := GenerateHash()

	room = configs.GetRouletteVariantMap(room.VariantType)

	room.ID = roomID
	room.Created = time.Now()
	room.CurrentRound = 1

	err = ro.Room.Create(room)
	if err != nil {
		writeErrorWithMsg(w, err)
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

func (ro *Room) RoomGetAllHandler(w http.ResponseWriter, r *http.Request) {

	rooms, err := ro.Room.GetAllRoom()
	if err != nil {
		writeErrorWithMsg(w, err)
		return
	}

	resp := struct {
		message string
		Data    []models.Room
	}{
		message: "Successfully created the room",
		Data:    rooms,
	}
	writeJSON(w, resp)
}

func (ro *Room) RoomGetHandler(w http.ResponseWriter, r *http.Request) {

	rMap := mux.Vars(r)

	roomid := rMap["id"]
	if len(roomid) == 0 {
		writeErrorWithMsg(w, errors.New("roomid must be provided"))
		return
	}

	room, err := ro.Room.GetRoom(roomid)
	if err != nil {
		writeErrorWithMsg(w, err)
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
