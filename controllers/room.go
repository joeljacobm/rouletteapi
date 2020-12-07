package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"rouletteapi/appconfigs"
	"rouletteapi/models"
	"rouletteapi/prometheus"
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

	var variant models.Variant
	var room models.Room
	err := json.NewDecoder(r.Body).Decode(&variant)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
		return
	}
	roomID := generateHash()

	variant = appconfigs.GetRouletteVariantMap(variant.VariantType)

	room.ID = roomID
	room.Created = time.Now()
	room.CurrentRound = 1
	room.RoomVariant = variant

	err = ro.Room.Create(room)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
		return
	}
	prometheus.RoomsTotal.Inc()

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
		if err == sql.ErrNoRows {
			writeErrorWithMsg(w, err, http.StatusNotFound)
			return
		}
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
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

func (ro *Room) RoomVariantHandler(w http.ResponseWriter, r *http.Request) {

	roomVariants := appconfigs.GetAllRouletteVariantMap()

	resp := struct {
		message string
		Data    []models.Variant
	}{
		message: "Successfully created the room",
		Data:    roomVariants,
	}
	writeJSON(w, resp)
}

func (ro *Room) RoomGetHandler(w http.ResponseWriter, r *http.Request) {

	rMap := mux.Vars(r)

	roomid := rMap["id"]
	if len(roomid) == 0 {
		writeErrorWithMsg(w, errors.New("roomid must be provided"), http.StatusNotFound)
		return
	}

	room, err := ro.Room.GetRoom(roomid)
	if err != nil {
		if err == sql.ErrNoRows {
			writeErrorWithMsg(w, err, http.StatusNotFound)
			return
		}
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
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
