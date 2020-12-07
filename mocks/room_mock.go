package mocks

import (
	"errors"
	"rouletteapi/models"
	"rouletteapi/postgres"
)

type RoomService struct {
	db postgres.DB
}

type BetService struct {
	db postgres.DB
}

func (rs RoomService) Create(room models.Room) error {
	return nil
}

func (rs RoomService) GetRoomCount(roomID string) (int, error) {
	if roomID == "" {
		return 0, errors.New("roomid cant be nil")
	}
	return 1, nil
}

func (rs RoomService) UpdateRound(roomID string) error {
	if roomID == "" {
		return errors.New("roomid cant be nil")
	}
	return nil
}

func (rs RoomService) GetRoom(roomid string) (models.Room, error) {
	return models.Room{ID: "5F9E314910",
		RoomVariant: models.Variant{VariantType: 1, VariantName: "sky_roulette"}}, nil
}

func (rs RoomService) GetAllRoom() ([]models.Room, error) {

	return []models.Room{{ID: "5F9E314910", RoomVariant: models.Variant{VariantType: 1, VariantName: "sky_roulette"}},
		{ID: "13D66E88E2", RoomVariant: models.Variant{VariantType: 2, VariantName: "superboost_roulett"}}}, nil

}
