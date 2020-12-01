package models

import (
	"rouletteapi/postgres"
	"time"
)

type PlayerService interface {
	Join(player Player) error
	GetPlayerCount(player Player) (int, error)
}

type playerService struct {
	db postgres.DB
}

type Player struct {
	ID          string    `json:"id"`
	RoomID      string    `json:"room_id"`
	ReadyStatus bool      `json:"ready_status"`
	BetsPlaced  []Bet     `json:"bets_placed"`
	DisplayName string    `json:"display_name"`
	InRoom      bool      `json:"in_room"`
	Created     time.Time `json:"created_at"`
}

func (ps playerService) Join(player Player) error {
	stmnt := ps.db.MustPrepare(`INSERT INTO player(created,player_id,room_id,in_room,ready_status,name) VALUES($1,$2,$3,$4,$5,$6);`)
	_, err := stmnt.Exec(player.Created, player.ID, player.RoomID, player.InRoom, player.ReadyStatus, player.DisplayName)
	return err
}

func (ps playerService) GetPlayerCount(player Player) (int, error) {
	var count int
	stmnt := ps.db.MustPrepare(`SELECT count(*) FROM player where player_id = $1 AND room_id = $2;`)
	row := stmnt.QueryRow(player.ID, player.RoomID)
	err := row.Scan(&count)
	return count, err
}
