package controllers

import (
	"rouletteapi/appconfigs"
	"rouletteapi/mocks"
	"rouletteapi/models"
	"rouletteapi/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	db      postgres.DB
	services = mocks.NewMockServices(db)
	roomC    = NewRoom(services.Room)
	playerC  = NewPlayer(services.Player, services.Room, services.Bet)
	betC     = NewBet(services.Player, services.Room, services.Bet)
)

func TestCalculateResult(t *testing.T) {

	testData1 := []models.Bet{
		{BetType: 1, Selection: 25, BetResult: models.Result{Number: 25, Colour: 2}},
		{BetType: 1, Selection: 25, BetResult: models.Result{Number: 6, Colour: 2}},
		{BetType: 2, Selection: 2, BetResult: models.Result{Number: 8, Colour: 2}},
		{BetType: 3, Selection: 1, BetResult: models.Result{Number: 8, Colour: 2}},
		{BetType: 3, Selection: 1, BetResult: models.Result{Number: 7, Colour: 2}},
		{BetType: 2, Selection: 1, BetResult: models.Result{Number: 8, Colour: 2}},
	}
	bets, err := calculateResult(testData1)
	if err != nil {
		t.Fatalf("CaluclateResult test failed with error %s", err)
	}

	bet1 := bets[0]
	assert.Equal(t, 1, bet1.Result)
	assert.Equal(t, "WIN", bet1.ResultText)

	bet2 := bets[1]
	assert.Equal(t, 0, bet2.Result)
	assert.Equal(t, "LOST", bet2.ResultText)

	bet3 := bets[2]
	assert.Equal(t, 1, bet3.Result)
	assert.Equal(t, "WIN", bet3.ResultText)

	bet4 := bets[3]
	assert.Equal(t, 1, bet4.Result)
	assert.Equal(t, "WIN", bet4.ResultText)

	bet5 := bets[4]
	assert.Equal(t, 0, bet5.Result)
	assert.Equal(t, "LOST", bet5.ResultText)

	bet6 := bets[5]
	assert.Equal(t, 0, bet6.Result)
	assert.Equal(t, "LOST", bet6.ResultText)

}

func TestCalculateLiability(t *testing.T) {

	appconfigs.LoadRouletteOddsMap("../appconfigs/oddsconfig.json")
	stake := 25.00
	liability := calculateLiability(1, stake)

	odds := 35.00 // straightup
	expectedLiability := odds * stake
	assert.Equal(t, expectedLiability, liability)
}
