package controllers

import (
	"encoding/json"
	"net/http"
	"rouletteapi/models"

	"github.com/gorilla/mux"
)

type Bet struct {
	Player models.PlayerService
	Room   models.RoomService
	Bet    models.BetService
	r      *mux.Router
}

func NewBet(player models.PlayerService, room models.RoomService, bet models.BetService, r *mux.Router) *Bet {
	return &Bet{
		Player: player,
		Room:   room,
		Bet:    bet,
		r:      r,
	}
}

func (b *Bet) BetResultHandler(w http.ResponseWriter, r *http.Request) {
	var bet models.Bet
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&bet)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}

	err = b.Bet.InsertResult(bet)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}
	err = b.Room.UpdateRound(bet.RoomID)
	if err != nil {
		writeErrorWithMsg(w, r, err)
		return
	}
	resp := struct {
		Message string
	}{
		Message: "Successfully inserted the result",
	}

	writeJSON(w, resp)

}
