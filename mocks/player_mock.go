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

	return nil, nil
}

func (ps PlayerService) Exit(player models.Player) error {

	return nil
}

func (ps PlayerService) CheckPlayerInRoomCount(id string) (int, error) {

	return 0, nil
}
