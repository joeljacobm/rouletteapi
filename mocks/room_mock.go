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
		Room:   mockRoomService{db: db},
		Player: mockPlayerService{db: db},
		Bet:    mockBetService{db: db},
	}
}

type mockRoomService struct {
	db postgres.DB
}

type mockPlayerService struct {
	db postgres.DB
}

type mockBetService struct {
	db postgres.DB
}

func (rs mockRoomService) Create(room models.Room) error {
	return nil
}

func (rs mockRoomService) GetRoomCount(RoomID string) (int, error) {
	return 1, nil
}

func (rs mockRoomService) UpdateRound(roomID string) error {

	return nil

}

func (ps mockPlayerService) Join(player models.Player) error {

	return nil
}

func (ps mockPlayerService) GetPlayer(id string, roomid string) (models.Player, error) {

	return models.Player{}, nil
}

// func (ps mockPlayerService) CheckInRoom(player models.Player) (bool, bool, error) {

// 	return true, false, nil
// }

func (ps mockPlayerService) GetReadyStatusForRound(roomid string) (int, error) {
	return 0, nil
}

func (bs mockBetService) PlaceBet(bet models.Bet, player models.Player) error {
	return nil
}

func (bs mockBetService) InsertResult(bet models.Bet) error {
	return nil
}

func (ps mockPlayerService) UpdateReadyStatus(player models.Player) error {

	return nil
}
