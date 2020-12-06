package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"rouletteapi/models"

	"github.com/gorilla/mux"
)

type Bet struct {
	Player models.PlayerService
	Room   models.RoomService
	Bet    models.BetService
}

func NewBet(player models.PlayerService, room models.RoomService, bet models.BetService) *Bet {
	return &Bet{
		Player: player,
		Room:   room,
		Bet:    bet,
	}
}

func (b *Bet) BetResultHandler(w http.ResponseWriter, r *http.Request) {
	var bet models.Bet
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&bet)
	if err != nil {
		writeErrorWithMsg(w, err)
		return
	}

	room, err := b.Room.GetRoom(bet.RoomID)
	if err != nil {
		writeErrorWithMsg(w, err)
		return
	}
	bet.RoundNo = room.CurrentRound

	err = b.Bet.InsertResult(bet)
	if err != nil {
		writeErrorWithMsg(w, err)
		return
	}
	err = b.Room.UpdateRound(bet.RoomID)
	if err != nil {
		writeErrorWithMsg(w, err)
		return
	}

	err = b.Player.UpdateReadyStatusFalse(bet.RoomID)
	if err != nil {
		writeErrorWithMsg(w, err)
		return
	}
	resp := struct {
		Message string
	}{
		Message: "Successfully inserted the result",
	}

	writeJSON(w, resp)

}

func (b *Bet) BetHandler(w http.ResponseWriter, r *http.Request) {

	rMap := mux.Vars(r)

	roomid := rMap["roomid"]
	if len(roomid) == 0 {
		writeErrorWithMsg(w, errors.New("roomid must be provided"))
		return
	}

	bets, err := b.Bet.GetBetForRoom(roomid)
	if err != nil {
		writeErrorWithMsg(w, err)
		return
	}
	resp := struct {
		Message string
		Data    []models.Bet
	}{
		Message: "Successfully retrieved the bets for the room",
		Data:    bets,
	}

	writeJSON(w, resp)

}
