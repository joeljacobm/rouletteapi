package models

import (
	"rouletteapi/postgres"
	"time"
)

// BetService Provides an interface for accessing the bet table
type BetService interface {
	PlaceBet(bet Bet, player Player) error
	InsertResult(bet Bet) error
	GetBet(id string, roomid string, roundno int) ([]Bet, error)
	GetBetForPlayer(roomid string) ([]Bet, error)
	GetBetForRoom(roomid string) ([]Bet, error)
}

type betService struct {
	db postgres.DB
}

// Bet implements the PlayerService Interface
type Bet struct {
	RoomID      string  `json:"room_id"`
	RoundNo     int     `json:"round_no"`
	BetType     int     `json:"bettype"`
	Stake       float64 `json:"stake"`
	Odds        float64 `json:"odds"`
	Liability   float64 `json:"liability"`
	Selection   int     `json:"selection"`
	Result      int     `json:"result"`
	ResultText  string  `json:"resultext"`
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
	stmnt := bs.db.MustPrepare(`INSERT INTO bet(created,round_no,bettype,selection,stake,liability,player_id,result,total_return,room_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`)
	_, err := stmnt.Exec(time.Now(), bet.RoundNo, bet.BetType, bet.Selection, bet.Stake, bet.Liability, player.ID, 0, 0, player.RoomID)
	return err
}

func (bs betService) InsertResult(bet Bet) error {
	stmnt := bs.db.MustPrepare(`INSERT INTO result(created,round_no,room_id,number,colour)VALUES($1,$2,$3,$4,$5);`)
	_, err := stmnt.Exec(time.Now(), bet.RoundNo, bet.RoomID, bet.BetResult.Number, bet.BetResult.Colour)
	return err
}

func (bs betService) GetBet(id string, roomid string, roundno int) ([]Bet, error) {
	var bet Bet
	var bets []Bet
	stmnt := bs.db.MustPrepare(`SELECT b.room_id,b.round_no, b.bettype,b.selection,b.stake,b.liability,r.number,r.colour FROM bet b LEFT JOIN result r ON b.room_id = r.room_id WHERE b.player_id = $1 AND r.room_id = $2 AND r.round_no = $3;`)
	rows, err := stmnt.Query(id, roomid, roundno)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&bet.RoomID, &bet.RoundNo, &bet.BetType, &bet.Selection, &bet.Stake, &bet.Liability, &bet.BetResult.Number, &bet.BetResult.Colour)
		if err != nil {
			return nil, err
		}
		bets = append(bets, bet)

	}
	return bets, err
}

func (bs betService) GetBetForRoom(roomid string) ([]Bet, error) {
	var bet Bet
	var bets []Bet
	stmnt := bs.db.MustPrepare(`SELECT b.room_id,b.round_no, b.bettype,b.selection,b.stake,b.liability,r.number,r.colour FROM bet b LEFT JOIN result r ON b.room_id = r.room_id WHERE r.room_id = $1;`)
	rows, err := stmnt.Query(roomid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&bet.RoomID, &bet.RoundNo, &bet.BetType, &bet.Selection, &bet.Stake, &bet.Liability, &bet.BetResult.Number, &bet.BetResult.Colour)
		if err != nil {
			return nil, err
		}
		bets = append(bets, bet)

	}

	return bets, err
}

func (bs betService) GetBetForPlayer(playerid string) ([]Bet, error) {
	var bet Bet
	var bets []Bet
	stmnt := bs.db.MustPrepare(`SELECT b.room_id,b.round_no, b.bettype,b.selection,b.stake,b.liability,r.number,r.colour FROM bet b LEFT JOIN result r ON b.room_id = r.room_id WHERE b.player_id = $1;`)
	rows, err := stmnt.Query(playerid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&bet.RoomID, &bet.RoundNo, &bet.BetType, &bet.Selection, &bet.Stake, &bet.Liability, &bet.BetResult.Number, &bet.BetResult.Colour)
		if err != nil {
			return nil, err
		}
		bets = append(bets, bet)

	}

	return bets, err
}
