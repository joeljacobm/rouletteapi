package models

import (
	"rouletteapi/postgres"
	"time"
)

type RoomService interface {
	Create(room Room) error
	GetRoomCount(roomID string) (int, error)
}

type roomService struct {
	db postgres.DB
}

type Room struct {
	ID           string    `json:"id"`
	VariantType  int       `json:"variant_type"` // Default is "skybet_roulette"
	VariantName  string    `json:"variant_name"`
	MaxPlayers   int       `json:"max_players"`
	NumOfRounds  int       `json:"max_rounds"`
	CurrentRound int       `json:"current_round"`
	Players      []Player  `json:"players"`
	Created      time.Time `json:"created_at"`
}

func (rs roomService) Create(room Room) error {
	stmnt := rs.db.MustPrepare(`INSERT INTO room(created,id,variant_type,variant_name,max_players,num_of_rounds,current_round) VALUES($1,$2,$3,$4,$5,$6,$7);`)
	_, err := stmnt.Exec(room.Created, room.ID, room.VariantType, room.VariantName, room.MaxPlayers, room.NumOfRounds, room.CurrentRound)
	return err
}

func (rs roomService) GetRoomCount(RoomID string) (int, error) {
	var count int
	stmnt := rs.db.MustPrepare(`SELECT count(*) FROM room where id = $1 AND num_of_rounds <> current_round;`)
	row := stmnt.QueryRow(RoomID)
	err := row.Scan(&count)
	return count, err
}
