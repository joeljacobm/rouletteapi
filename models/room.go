package models

import (
	"rouletteapi/postgres"
	"time"
)

// RoomService Provides an interface for accessing the room table
type RoomService interface {
	Create(room Room) error
	GetRoomCount(roomID string) (int, error)
	UpdateRound(roomID string) error
	GetAllRoom() ([]Room, error)
	GetRoom(roomid string) (Room, error)
}

type roomValidator struct {
	RoomService
}

type roomService struct {
	db postgres.DB
}

// Room implements the RoomService interface
type Room struct {
	ID           string    `json:"id"`
	RoomVariant  Variant   `json:"variant"`
	CurrentRound int       `json:"current_round"`
	Players      []Player  `json:"players"`
	Created      time.Time `json:"created_at"`
}

// Variant is the roulette variant type
type Variant struct {
	VariantType int    `json:"variant_type"` // Default is "skybet_roulette"
	VariantName string `json:"variant_name"`
	MaxPlayers  int    `json:"max_players"`
	NumOfRounds int    `json:"max_rounds"`
}

func (rs roomService) Create(room Room) error {
	stmnt := rs.db.MustPrepare(`INSERT INTO room(created,id,variant_type,variant_name,max_players,num_of_rounds,current_round) VALUES($1,$2,$3,$4,$5,$6,$7);`)
	_, err := stmnt.Exec(room.Created, room.ID, room.RoomVariant.VariantType, room.RoomVariant.VariantName, room.RoomVariant.MaxPlayers, room.RoomVariant.NumOfRounds, room.CurrentRound)
	return err
}

func (rs roomService) GetRoomCount(RoomID string) (int, error) {
	var count int
	stmnt := rs.db.MustPrepare(`SELECT count(*) FROM room where id = $1 AND num_of_rounds <> current_round;`)
	row := stmnt.QueryRow(RoomID)
	err := row.Scan(&count)
	return count, err
}

func (rs roomService) UpdateRound(roomID string) error {
	stmnt := rs.db.MustPrepare(`UPDATE room SET current_round = current_round+1 WHERE id = $1;`)
	_, err := stmnt.Exec(roomID)
	return err

}

func (rs roomService) GetRoom(roomid string) (Room, error) {
	var room Room
	stmnt := rs.db.MustPrepare(`SELECT id,variant_type,variant_name,max_players,num_of_rounds,current_round FROM room WHERE id = $1;`)
	row := stmnt.QueryRow(roomid)
	err := row.Scan(&room.ID, &room.RoomVariant.VariantType, &room.RoomVariant.VariantName, &room.RoomVariant.MaxPlayers, &room.RoomVariant.NumOfRounds, &room.CurrentRound)
	return room, err

}

func (rs roomService) GetAllRoom() ([]Room, error) {
	var room Room
	var rooms []Room
	stmnt := rs.db.MustPrepare(`SELECT id,variant_type,variant_name,max_players,num_of_rounds,current_round FROM room WHERE current_round <> num_of_rounds;`)
	rows, err := stmnt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&room.ID, &room.RoomVariant.VariantType, &room.RoomVariant.VariantName, &room.RoomVariant.MaxPlayers, &room.RoomVariant.NumOfRounds, &room.CurrentRound)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)

	}

	return rooms, err

}
