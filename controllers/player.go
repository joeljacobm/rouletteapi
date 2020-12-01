package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rouletteapi/configs"
	"rouletteapi/models"
	"time"

	"github.com/gorilla/mux"
)

type Player struct {
	Player models.PlayerService
	Room   models.RoomService
	Bet    models.BetService
	r      *mux.Router
}

func NewPlayer(player models.PlayerService, room models.RoomService, bet models.BetService, r *mux.Router) *Player {
	return &Player{
		Player: player,
		Room:   room,
		Bet:    bet,
		r:      r,
	}
}

func (pl *Player) PlayerGetBetHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	playerid, ok := query["player_id"]
	if !ok || len(playerid) == 0 {
		writeErrorWithMsg(w, r, errors.New(""))
	}

	roomid, ok := query["room_id"]
	if !ok || len(roomid) == 0 {
		writeErrorWithMsg(w, r, errors.New(""))
	}

	round, ok := query["round"]
	if !ok || len(round) == 0 {
		writeErrorWithMsg(w, r, errors.New(""))
	}

	player, err := pl.Player.GetPlayer(playerid[0], roomid[0])
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}
	fmt.Println(player)

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

	// ok, err = checkIfPlayerIsInRoom(pl, player)

	// if err != nil || !ok {
	// 	writeErrorWithMsg(w, r, err)
	// 	return
	// }

	player.Created = time.Now()
	player.InRoom = true

	pl.Player.Join(player)
	if err != nil {
		fmt.Println(err)

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
		return false, errors.New("Room not present")
	}
	return true, nil
}

func checkIfPlayerIsInRoom(pl *Player, p models.Player) (bool, error) {

	player, err := pl.Player.GetPlayer(p.ID, p.RoomID)
	if err != nil {
		return false, err
	}

	if player.ID != "" {
		return false, errors.New("Player is already in the room")
	}
	return true, nil
}

func (pl *Player) PlayerReadyHandler(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}

	err = pl.Player.UpdateReadyStatus(player)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}
	var count int
	count, err = pl.Player.GetReadyStatusForRound(player.RoomID)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}

	var resp struct{ Message string }
	if count != 0 {
		resp = struct {
			Message string
		}{
			Message: "Wait for the all the players to be ready",
		}
	} else {
		resp = struct {
			Message string
		}{
			Message: "All the players are ready. Lets spin the ball",
		}

	}

	writeJSON(w, resp)

}

func (pl *Player) PlayerBetHandler(w http.ResponseWriter, r *http.Request) {

	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}

	p, err := pl.Player.GetPlayer(player.ID, player.RoomID)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}

	if !p.InRoom {
		writeErrorWithMsg(w, r, errors.New("Player not in the room"))
		return
	}

	if p.ReadyStatus {
		writeErrorWithMsg(w, r, errors.New("Bet not accepted since player is already ready"))
		return
	}

	for _, bet := range player.BetsPlaced {
		bet.Liability = calculateLiability(bet.BetType, bet.Stake)

		err := pl.Bet.PlaceBet(bet, player)
		if err != nil {
			writeErrorWithMsg(w, r, err)
			return
		}
	}

	resp := struct {
		Message string
	}{
		Message: "Succefully placed the bet",
	}
	writeJSON(w, resp)

}

func calculateLiability(bettype int, stake float64) float64 {
	odds := configs.GetRouletteOddsMap(bettype)
	return odds * stake
}
