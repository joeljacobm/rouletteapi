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

type RoomService struct {
	db postgres.DB
}

type PlayerService struct {
	db postgres.DB
}

type PlayerValidator struct {
	models.PlayerService
}

type BetService struct {
	db postgres.DB
}

func (rs RoomService) Create(room models.Room) error {
	return nil
}

func (rs RoomService) GetRoomCount(RoomID string) (int, error) {
	return 1, nil
}

func (rs RoomService) UpdateRound(roomID string) error {

	return nil

}

// func (pv PlayerValidator) Join(player models.Player) error {

// 	if len(player.ID) < 10 {
// 		return errors.New("Player ID must be 10 characters in length")
// 	}
// 	return pv.Join(player)
// }

func (ps PlayerService) Join(player models.Player) error {

	return nil
}

func (ps PlayerService) GetPlayer(id string, roomid string) (models.Player, error) {

	return models.Player{}, nil
}

// func (ps mockPlayerService) CheckInRoom(player models.Player) (bool, bool, error) {

// 	return true, false, nil
// }

func (ps PlayerService) GetReadyStatusForRound(roomid string) (int, error) {
	return 0, nil
}

func (bs BetService) PlaceBet(bet models.Bet, player models.Player) error {
	return nil
}

func (bs BetService) InsertResult(bet models.Bet) error {
	return nil
}

func (ps PlayerService) UpdateReadyStatusTrue(player models.Player) error {

	return nil
}


func (ps PlayerService) UpdateReadyStatusFalse(roomid string) error {
	
	return nil
}

func (bs BetService) GetBet(id string, roomid string, roundno int) ([]models.Bet, error) {

	return []models.Bet{{RoomID: "100D2BF54B", RoundNo: 1, BetType: 1, Stake: 1.5, Liability: 52.5, Selection: 24,
		BetResult: models.Result{Number: 24, Colour: 2}}, {RoomID: "100D2BF54B", RoundNo: 1, BetType: 2, Stake: 1.5, Liability: 52.5, Selection: 1,
		BetResult: models.Result{Number: 24, Colour: 2}}}, nil
}
