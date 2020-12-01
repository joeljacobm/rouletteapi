package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"rouletteapi/models"
	"time"

	"github.com/gorilla/mux"
)

type Player struct {
	Player models.PlayerService
	Room   models.RoomService
	r      *mux.Router
}

func NewPlayer(player models.PlayerService, room models.RoomService, r *mux.Router) *Player {
	return &Player{
		Player: player,
		Room:   room,
		r:      r,
	}
}

func (pl *Player) PlayerJoinHandler(w http.ResponseWriter, r *http.Request) {

	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}

	if player.DisplayName == "" {
		dName := "Display" + GenerateHash()
		player.DisplayName = dName
	}

	ok, err := checkIfRoomIsActive(pl, player)
	if err != nil || !ok {
		writeErrorWithMsg(w, r, err)
		return
	}

	ok, err = checkIfPlayerIsInRoom(pl, player)
	if err != nil || !ok {
		writeErrorWithMsg(w, r, err)
		return
	}

	player.Created = time.Now()
	player.InRoom = true

	pl.Player.Join(player)
	if err != nil {
		writeErrorWithMsg(w, r, errors.New("Cannot join the room"))
		return
	}
	resp := struct {
		Message string
		Data    models.Player
	}{
		Message: "Successfully joined the room",
		Data:    player,
	}
	writeJSON(w, resp)
}

func checkIfRoomIsActive(pl *Player, p models.Player) (bool, error) {

	var count int
	count, err := pl.Room.GetRoomCount(p.RoomID)
	if err != nil {
		return false, errors.New("Room does not exist")
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

func checkIfPlayerIsInRoom(pl *Player, p models.Player) (bool, error) {

	var count int
	count, err := pl.Player.GetPlayerCount(p)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, errors.New("Player is already in the room")
	}
	return true, nil
}
