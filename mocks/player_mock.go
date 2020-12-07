package mocks

import (
	"errors"
	"rouletteapi/models"
	"rouletteapi/postgres"
)

type PlayerService struct {
	db postgres.DB
}

type PlayerValidator struct {
	models.PlayerService
}

func (pv PlayerValidator) Join(player models.Player) error {

	if len(player.ID) < 10 {
		return errors.New("Player ID must be 10 characters in length")
	}
	return pv.PlayerService.Join(player)
}

func (ps PlayerService) Join(player models.Player) error {

	return nil
}

func (ps PlayerService) GetPlayer(id string, roomid string) (models.Player, error) {

	if id == "6EB9E64G56" {
		return models.Player{InRoom: true}, nil
	}
	return models.Player{}, nil
}

func (ps PlayerService) GetReadyStatusForRound(roomid string) (int, error) {
	return 0, nil
}

func (ps PlayerService) UpdateReadyStatusFalse(roomid string) error {

	return nil
}

func (ps PlayerService) UpdateReadyStatusTrue(player models.Player) error {

	return nil
}

func (ps PlayerService) GetAllPlayers() ([]models.Player, error) {

	return []models.Player{{RoomID: "100D2BF54B", ID: "6EB8E64F56", InRoom: true}, {RoomID: "200F2BF54B", ID: "7EB9E64F56", InRoom: false}}, nil
}

func (ps PlayerService) Exit(player models.Player) error {

	if player.ID != "5EB8E64F56" {
		return errors.New("Wrong player id")
	}
	return nil
}

func (ps PlayerService) CheckPlayerInRoomCount(id string) (int, error) {

	return 0, nil
}
