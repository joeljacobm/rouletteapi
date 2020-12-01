package mocks

import (
	"rouletteapi/models"
	"rouletteapi/postgres"
)

type mockServices struct {
	Room   models.RoomService
	Player models.PlayerService
}

func NewMockServices(db postgres.DB) mockServices {
	return mockServices{
		Room:   mockRoomService{db: db},
		Player: mockPlayerService{db: db},
	}
}

type mockRoomService struct {
	db postgres.DB
}

type mockPlayerService struct {
	db postgres.DB
}

func (rs mockRoomService) Create(room models.Room) error {
	return nil
}

func (rs mockRoomService) GetRoomCount(RoomID string) (int, error) {
	return 1, nil
}

func (ps mockPlayerService) Join(player models.Player) error {

	return nil
}

func (ps mockPlayerService) GetPlayerCount(player models.Player) (int, error) {

	return 0, nil
}
