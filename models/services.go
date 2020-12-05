package models

import (
	"rouletteapi/postgres"
)

type Services struct {
	Room   RoomService
	Player PlayerService
	Bet    BetService
}

func NewServices(db postgres.DB) Services {
	return Services{
		Room:   roomService{db},
		Player: playerValidator{PlayerService: playerService{db: db}},
		Bet:    betService{db: db},
	}
}

// func NewRoomService(db postgres.DB) RoomService {
// 	return
// }

// func NewPlayerService(db postgres.DB) PlayerService {
// 	return
// }
