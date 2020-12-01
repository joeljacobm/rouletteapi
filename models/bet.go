package models

import (
	"rouletteapi/postgres"
	"time"
)

type BetService interface {
	PlaceBet(bet Bet, player Player) error
	InsertResult(bet Bet) error
}

type betService struct {
	db postgres.DB
}
type Bet struct {
	RoomID      string  `json:"room_id"`
	RoundNo     int     `json:"round_no"`
	BetType     int     `json:"bettype"`
	Stake       float64 `json:"stake"`
	Odds        float64 `json:"odds"`
	Liability   float64 `json:"liability"`
	Selection   int     `json:"selection"`
	Result      int     `json:"result"`
	TotalReturn float64 `json:"total_return"`
	BetResult   Result  `json:"bet_result"`
}

type Result struct {
	Number  int `json:"number"`
	Colour  int `json:"colour"`
	OddEven int `json:"oddeven"`
}

type BetType struct {
	BetName     string  `json:"betname"`
	OddsDecimal float64 `json:"oddsdecimal"`
}

func (bs betService) PlaceBet(bet Bet, player Player) error {
	stmnt := bs.db.MustPrepare(`INSERT INTO bet(created,round_no,bettype,stake,liability,player_id,result,total_return,room_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);`)
	_, err := stmnt.Exec(time.Now(), bet.RoundNo, bet.BetType, bet.Stake, bet.Liability, player.ID, 0, 0, player.RoomID)
	return err
}

func (bs betService) InsertResult(bet Bet) error {
	stmnt := bs.db.MustPrepare(`INSERT INTO result(created,round_no,room_id,number,colour)VALUES($1,$2,$3,$4,$5);`)
	_, err := stmnt.Exec(time.Now(), bet.RoundNo,bet.RoomID, bet.BetResult.Number, bet.BetResult.Colour)
	return err
}
