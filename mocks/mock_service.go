package mocks

import (
	"rouletteapi/models"
	"rouletteapi/postgres"
)

type mockServices struct {
	Room   models.RoomService
	Player models.PlayerService
	Bet    models.BetService
}

func NewMockServices(db postgres.DB) mockServices {
	return mockServices{
		Room:   RoomService{db: db},
		Player: PlayerValidator{PlayerService: PlayerService{db: db}},
		Bet:    BetService{db: db},
	}
}
