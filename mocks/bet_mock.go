package mocks

import (
	"errors"
	"rouletteapi/models"
)

func (bs BetService) PlaceBet(bet models.Bet, player models.Player) error {
	return nil
}

func (bs BetService) InsertResult(bet models.Bet) error {
	return nil
}

func (bs BetService) GetBet(id string, roomid string, roundno int) ([]models.Bet, error) {

	if id == "5EB8E64F56" {
		return []models.Bet{{RoomID: "100D2BF54B", RoundNo: 1, BetType: 1, Stake: 1.5, Liability: 52.5, Selection: 24,
			BetResult: models.Result{Number: 24, Colour: 2}}, {RoomID: "100D2BF54B", RoundNo: 1, BetType: 2, Stake: 1.5, Liability: 52.5, Selection: 1,
			BetResult: models.Result{Number: 24, Colour: 2}}}, nil
	}

	if id == "6EB8E64F56" {
		return []models.Bet{{RoomID: "100D2BF54B", RoundNo: 1, BetType: 1, Stake: 1.5, Liability: 52.5, Selection: 24,
			BetResult: models.Result{Number: 24, Colour: 2}}, {RoomID: "100D2BF54B", RoundNo: 1, BetType: 2, Stake: 1.5, Liability: 52.5, Selection: 1,
			BetResult: models.Result{Number: 24, Colour: 2}}}, nil
	}

	if id == "7EB9E64F56" {
		return []models.Bet{{RoomID: "200F2BF54B", RoundNo: 1, BetType: 1, Stake: 1.5, Liability: 52.5, Selection: 24,
			BetResult: models.Result{Number: 24, Colour: 2}}, {RoomID: "200F2BF54B", RoundNo: 1, BetType: 3, Stake: 1.5, Liability: 52.5, Selection: 1,
			BetResult: models.Result{Number: 24, Colour: 2}}}, nil
	}
	return nil, nil
}

func (bs BetService) GetBetForPlayer(playerid string) ([]models.Bet, error) {

	if playerid == "6EB8E64F56" {
		return []models.Bet{
			{RoomID: "100D2BF54B", RoundNo: 1, BetType: 1, Stake: 1.5, Selection: 24},
			{RoomID: "100D2BF54B", RoundNo: 1, BetType: 2, Stake: 1.5, Selection: 1},
		}, nil
	}

	if playerid == "7EB9E64F56" {
		return []models.Bet{
			{RoomID: "200F2BF54B", RoundNo: 1, BetType: 1, Stake: 1.5, Selection: 24},
			{RoomID: "200F2BF54B", RoundNo: 1, BetType: 2, Stake: 1.5, Selection: 1},
		}, nil
	}
	return nil, nil
}

func (bs BetService) GetBetForRoom(roomid string) ([]models.Bet, error) {

	if roomid == "100D2BF54B" {
		return []models.Bet{{RoomID: "100D2BF54B", RoundNo: 1, BetType: 1, Stake: 1.5, Selection: 24},
			{RoomID: "100D2BF54B", RoundNo: 1, BetType: 2, Stake: 1.5, Selection: 1}}, nil
	}
	return nil, errors.New("Invalid roomid")
}
