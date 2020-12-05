package controllers

import (
	"encoding/json"
	"net/http"
	"rouletteapi/models"
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
