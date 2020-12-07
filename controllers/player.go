package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"rouletteapi/models"
	"rouletteapi/prometheus"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Player struct {
	Player models.PlayerService
	Room   models.RoomService
	Bet    models.BetService
}

func NewPlayer(player models.PlayerService, room models.RoomService, bet models.BetService) *Player {
	return &Player{
		Player: player,
		Room:   room,
		Bet:    bet,
	}
}

func (pl *Player) PlayerGetResultHandler(w http.ResponseWriter, r *http.Request) {

	pMap := mux.Vars(r)
	playerid := pMap["playerid"]
	if len(playerid) == 0 {
		writeErrorWithMsg(w, errors.New("playerid must be provided"), http.StatusBadRequest)
		return
	}

	roomid := pMap["roomid"]
	if len(roomid) == 0 {
		writeErrorWithMsg(w, errors.New("roomid must be provided"), http.StatusBadRequest)
		return
	}

	roundno := pMap["roundno"]
	roundInt, err := strconv.Atoi(roundno)
	if roundInt == 0 || err != nil {
		writeErrorWithMsg(w, errors.New("round number must be provided"), http.StatusBadRequest)
		return
	}

	bet, err := pl.Bet.GetBet(playerid, roomid, roundInt)
	if err != nil {
		checkError(w, err)
		return
	}

	bets, err := calculateResult(bet)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
		return
	}

	resp := struct {
		Message string
		Data    []models.Bet
	}{
		Message: "Successfully retrieved the result",
		Data:    bets,
	}
	writeJSON(w, resp)

}

func (pl *Player) PlayerJoinHandler(w http.ResponseWriter, r *http.Request) {

	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
		return
	}

	if player.DisplayName == "" {
		dName := "Roulette-" + generateHash() + "-Master"
		player.DisplayName = dName
	}

	ok, err := checkIfRoomIsActive(pl, player)
	if err != nil || !ok {
		writeErrorWithMsg(w, err, http.StatusNotFound)
		return
	}

	ok, err = checkIfPlayerIsInRoom(pl, player)
	if err != nil || !ok {
		writeErrorWithMsg(w, err, http.StatusNotFound)
		return
	}

	player.Created = time.Now()
	player.InRoom = true

	err = pl.Player.Join(player)
	if err != nil {
		writeErrorWithMsg(w, errors.New("Cannot join the room"), http.StatusNotFound)
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
		return false, errors.New("Room does not exist")
	}
	return true, nil
}

func checkIfPlayerIsInRoom(pl *Player, p models.Player) (bool, error) {

	count, err := pl.Player.CheckPlayerInRoomCount(p.ID)
	if count > 0 || err != nil {
		return false, errors.New("Player is already in another room. Please exit the room first")

	}

	player, err := pl.Player.GetPlayer(p.ID, p.RoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	if player.InRoom {
		return false, errors.New("Player is already in the room")
	}

	return true, nil
}

func (pl *Player) PlayerReadyHandler(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
		return
	}

	err = pl.Player.UpdateReadyStatusTrue(player)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
		return
	}
	var count int
	count, err = pl.Player.GetReadyStatusForRound(player.RoomID)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
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
			Message: "All the players are ready. Ready to spin",
		}

	}

	writeJSON(w, resp)

}

func (pl *Player) PlayerBetHandler(w http.ResponseWriter, r *http.Request) {

	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
		return
	}

	p, err := pl.Player.GetPlayer(player.ID, player.RoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeErrorWithMsg(w, errors.New("Player is not in the room"), http.StatusNotFound)
		}
		return
	}

	if !p.InRoom {
		writeErrorWithMsg(w, errors.New("Player is not in the room"), http.StatusNotFound)
		return
	}

	if p.ReadyStatus {
		writeErrorWithMsg(w, errors.New("Bet not accepted since player is already ready for the current spin"), http.StatusNotFound)
		return
	}

	room, err := pl.Room.GetRoom(player.RoomID)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
		return
	}

	for _, bet := range player.BetsPlaced {
		bet.Liability = calculateLiability(bet.BetType, bet.Stake)

		bet.RoundNo = room.CurrentRound
		err := pl.Bet.PlaceBet(bet, player)
		if err != nil {
			writeErrorWithMsg(w, err, http.StatusInternalServerError)
			return
		}
	}

	resp := struct {
		Message string
	}{
		Message: "Succefully placed the bet",
	}
	prometheus.BetsTotal.Inc()
	writeJSON(w, resp)

}

func (pl *Player) PlayerHandler(w http.ResponseWriter, r *http.Request) {

	players, err := pl.Player.GetAllPlayers()
	if err != nil {
		checkError(w, err)
		return
	}

	for k, p := range players {
		playerBets, err := pl.Bet.GetBetForPlayer(p.ID)
		if err != nil {
			checkError(w, err)
			return
		}
		var round int
		for _, b := range playerBets {

			if b.RoundNo == round {
				continue
			}

			round = b.RoundNo
			bet, err := pl.Bet.GetBet(p.ID, p.RoomID, b.RoundNo)
			if err != nil {
				checkError(w, err)
				return
			}
			bets, err := calculateResult(bet)
			if err != nil {
				writeErrorWithMsg(w, err, http.StatusInternalServerError)
				return
			}
			players[k].BetsPlaced = append(players[k].BetsPlaced, bets...)

		}

	}

	resp := struct {
		Message string
		Data    []models.Player
	}{
		Message: "Succefully retrieved the player details",
		Data:    players,
	}
	writeJSON(w, resp)

}

func (pl *Player) PlayerExitHandler(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusInternalServerError)
		return
	}

	err = pl.Player.Exit(player)
	if err != nil {
		writeErrorWithMsg(w, err, http.StatusNotFound)
		return
	}

	var resp struct{ Message string }
	resp = struct {
		Message string
	}{
		Message: "Successfully exited from the room",
	}

	writeJSON(w, resp)

}
