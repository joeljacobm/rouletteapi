package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rouletteapi/configs"
	"rouletteapi/mocks"
	"rouletteapi/models"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestPlayerJoinHandler(t *testing.T) {
	testData := struct {
		RoomID string `json:"room_id"`
		ID     string `json:"id"`
	}{
		RoomID: "FA399BB381",
		ID:     "5EB8E63F56",
	}

	services := mocks.NewMockServices(db)
	playerC := NewPlayer(services.Player, services.Room, services.Bet)
	b, _ := json.Marshal(testData)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/player/join", bytes.NewReader(b))
	r.Header.Add("Content-Type", "application/json")
	playerC.PlayerJoinHandler(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	expected := struct {
		message string
		Data    models.Player
	}{}
	json.Unmarshal(resp, &expected)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, true, expected.Data.InRoom)
	assert.Equal(t, false, expected.Data.ReadyStatus)
	assert.Equal(t, "FA399BB381", expected.Data.RoomID)
	assert.Equal(t, "5EB8E63F56", expected.Data.ID)
}

func TestPlayerBetHandler(t *testing.T) {
	testData := models.Player{
		RoomID: "100D2BF54B",
		ID:     "5EB8E64F56",
		BetsPlaced: []models.Bet{
			{RoundNo: 1, BetType: 1, Stake: 1.5, Selection: 24},
			{RoundNo: 1, BetType: 2, Stake: 1.5, Selection: 1},
		},
	}

	services := mocks.NewMockServices(db)
	playerC := NewPlayer(services.Player, services.Room, services.Bet)
	b, _ := json.Marshal(testData)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/player/join", bytes.NewReader(b))
	r.Header.Add("Content-Type", "application/json")
	playerC.PlayerJoinHandler(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	expected := struct {
		message string
		Data    models.Player
	}{}
	json.Unmarshal(resp, &expected)

	assert.Equal(t, 200, w.Result().StatusCode)
}

func TestPlayerGetBetHandler(t *testing.T) {

	configs.LoadRouletteOddsMap("../configs/oddsconfig.json")
	services := mocks.NewMockServices(db)
	playerC := NewPlayer(services.Player, services.Room, services.Bet)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/player/result", nil)

	vars := map[string]string{
		"playerid": "5EB8E64F56",
		"roomid":   "100D2BF54B",
		"roundno":  "1",
	}

	r = mux.SetURLVars(r, vars)
	r.Header.Add("Content-Type", "application/json")
	playerC.PlayerGetBetHandler(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	expected := struct {
		message string
		Data    []models.Bet
	}{}

	json.Unmarshal(resp, &expected)

	bet1 := expected.Data[0]
	bet2 := expected.Data[1]

	assert.Equal(t, 54.00, bet1.TotalReturn)
	assert.Equal(t, 35.00, bet1.Odds)
	assert.Equal(t, 1, bet1.Result)
	assert.Equal(t, "WIN", bet1.ResultText)
	assert.Equal(t, 24, bet1.BetResult.Number)

	assert.Equal(t, "LOST", bet2.ResultText)
	assert.Equal(t, 0, bet2.Result)
	assert.Equal(t, 0.00, bet2.TotalReturn)
	assert.Equal(t, 24, bet2.BetResult.Number)
	assert.Equal(t, 2, bet2.BetResult.Colour)

	assert.Equal(t, 200, w.Result().StatusCode)

}
