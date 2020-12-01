package models

import (
	"rouletteapi/postgres"
)

type services struct {
	Room   RoomService
	Player PlayerService
	Bet    BetService
}

func NewServices(db postgres.DB) services {
	return services{
		Room: roomService{db: db},
		Player:playerService{db: db},
		Bet: betService{db: db},
	}
}

// func NewRoomService(db postgres.DB) RoomService {
// 	return 
// }

// func NewPlayerService(db postgres.DB) PlayerService {
// 	return 
// }



