package models

import (
	"errors"
	"rouletteapi/postgres"
	"time"
)

type PlayerService interface {
	Join(player Player) error
	GetPlayer(id string, roomid string) (Player, error)
	// CheckInRoom(player Player) (bool, bool, error)
	GetReadyStatusForRound(roomid string) (int, error)
	UpdateReadyStatusTrue(player Player) error
	UpdateReadyStatusFalse(roomid string) error
}

type playerValidator struct {
	PlayerService
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

func (pv playerValidator) Join(player Player) error {

	if len(player.ID) < 10 {
		return errors.New("Player ID must be atleast 10 characters long")
	}
	return pv.PlayerService.Join(player)
}

func (ps playerService) Join(player Player) error {
	stmnt := ps.db.MustPrepare(`INSERT INTO player(created,player_id,room_id,in_room,ready_status,name) VALUES($1,$2,$3,$4,$5,$6);`)
	_, err := stmnt.Exec(player.Created, player.ID, player.RoomID, player.InRoom, player.ReadyStatus, player.DisplayName)
	return err
}

func (ps playerService) GetPlayer(id string, roomid string) (Player, error) {
	var player Player
	stmnt := ps.db.MustPrepare(`SELECT * FROM player WHERE player_id = $1 AND room_id = $2 limit 1;`)
	row := stmnt.QueryRow(id, roomid)
	err := row.Scan(&player.Created, &player.ID, &player.RoomID, &player.InRoom, &player.ReadyStatus, &player.DisplayName)
	return player, err
}

// func (ps playerService) CheckInRoom(player Player) (bool, bool, error) {
// 	var inroom, readystatus bool
// 	stmnt := ps.db.MustPrepare(`SELECT in_room,ready_status FROM player WHERE player_id = $1 AND room_id = $2;`)
// 	row := stmnt.QueryRow(player.ID, player.RoomID)
// 	err := row.Scan(&inroom, &readystatus)
// 	return inroom, readystatus, err
// }

func (ps playerService) GetReadyStatusForRound(roomid string) (int, error) {
	var count int
	stmnt := ps.db.MustPrepare(`SELECT count(*) FROM player WHERE room_id = $1 AND ready_status = false;`)
	row := stmnt.QueryRow(roomid)
	err := row.Scan(&count)
	return count, err
}

func (ps playerService) UpdateReadyStatusTrue(player Player) error {
	stmnt := ps.db.MustPrepare(`UPDATE player set ready_status = true WHERE player_id = $1 AND room_id = $2;`)
	_, err := stmnt.Exec(player.ID, player.RoomID)
	return err
}

func (ps playerService) UpdateReadyStatusFalse(roomid string) error {
	stmnt := ps.db.MustPrepare(`UPDATE player set ready_status = false WHERE room_id = $1;`)
	_, err := stmnt.Exec(roomid)
	return err
}
