package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rouletteapi/appconfigs"
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
		ID:     "6EB9E64G56",
		BetsPlaced: []models.Bet{
			{RoundNo: 1, BetType: 1, Stake: 1.5, Selection: 24},
			{RoundNo: 1, BetType: 2, Stake: 1.5, Selection: 1},
		},
	}

	b, _ := json.Marshal(testData)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/player/join", bytes.NewReader(b))
	r.Header.Add("Content-Type", "application/json")
	playerC.PlayerBetHandler(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	expected := struct {
		message string
		Data    models.Player
	}{}
	json.Unmarshal(resp, &expected)

	assert.Equal(t, 200, w.Result().StatusCode)
}

func TestPlayerGetResultHandler(t *testing.T) {

	appconfigs.LoadRouletteOddsMap("../appconfigs/oddsconfig.json")

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/player/result", nil)

	vars := map[string]string{
		"playerid": "5EB8E64F56",
		"roomid":   "100D2BF54B",
		"roundno":  "1",
	}

	r = mux.SetURLVars(r, vars)
	r.Header.Add("Content-Type", "application/json")
	playerC.PlayerGetResultHandler(w, r)

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

func TestPlayerReadyHandler(t *testing.T) {

	testData := models.Player{
		RoomID: "100D2BF54B",
		ID:     "5EB8E64F56",
	}
	b, _ := json.Marshal(testData)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/player/ready", bytes.NewReader(b))

	r.Header.Add("Content-Type", "application/json")
	playerC.PlayerReadyHandler(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	expected := struct {
		message string
		Data    models.Player
	}{}
	json.Unmarshal(resp, &expected)
	assert.Equal(t, 200, w.Result().StatusCode)

}

func TestPlayerHandler(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/player", nil)

	r.Header.Add("Content-Type", "application/json")
	playerC.PlayerHandler(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	expectedData := []models.Player{
		{RoomID: "100D2BF54B", ID: "6EB8E64F56",
			BetsPlaced: []models.Bet{
				{RoundNo: 1, BetType: 1, Stake: 1.5, Selection: 24, ResultText: "WIN", TotalReturn: 54},
				{RoundNo: 1, BetType: 2, Stake: 1.5, Selection: 1, ResultText: "LOST", TotalReturn: 0},
			}, InRoom: true}, {RoomID: "200F2BF54B", ID: "7EB9E64F56",
			BetsPlaced: []models.Bet{
				{RoundNo: 1, BetType: 1, Stake: 1.5, Selection: 24, ResultText: "WIN", TotalReturn: 54},
				{RoundNo: 1, BetType: 3, Stake: 1.5, Selection: 1, ResultText: "WIN", TotalReturn: 3},
			}, InRoom: false},
	}
	actual := struct {
		message string
		Data    []models.Player
	}{}
	json.Unmarshal(resp, &actual)

	player1 := actual.Data[0]

	player2 := actual.Data[1]
	player1bet1 := player1.BetsPlaced[0]
	player1bet2 := player1.BetsPlaced[1]
	player2bet1 := player2.BetsPlaced[0]
	player2bet2 := player2.BetsPlaced[1]

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, expectedData[0].ID, player1.ID)
	assert.Equal(t, expectedData[0].RoomID, player1.RoomID)
	assert.Equal(t, expectedData[0].ReadyStatus, player1.ReadyStatus)
	assert.Equal(t, expectedData[0].BetsPlaced[0].ResultText, player1bet1.ResultText)
	assert.Equal(t, expectedData[0].BetsPlaced[1].ResultText, player1bet2.ResultText)
	assert.Equal(t, expectedData[0].BetsPlaced[0].TotalReturn, player1bet1.TotalReturn)
	assert.Equal(t, expectedData[0].BetsPlaced[1].TotalReturn, player1bet2.TotalReturn)

	assert.Equal(t, expectedData[1].ID, player2.ID)
	assert.Equal(t, expectedData[1].RoomID, player2.RoomID)
	assert.Equal(t, expectedData[1].ReadyStatus, player2.ReadyStatus)
	assert.Equal(t, expectedData[1].BetsPlaced[0].ResultText, player2bet1.ResultText)
	assert.Equal(t, expectedData[1].BetsPlaced[1].ResultText, player2bet2.ResultText)
	assert.Equal(t, expectedData[1].BetsPlaced[0].TotalReturn, player2bet1.TotalReturn)
	assert.Equal(t, expectedData[1].BetsPlaced[1].TotalReturn, player2bet2.TotalReturn)

}

func TestPlayerExitHandler(t *testing.T) {

	testData := models.Player{
		ID:     "5EB8E64F56",
	}
	b, _ := json.Marshal(testData)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/player/exit", bytes.NewReader(b))

	r.Header.Add("Content-Type", "application/json")
	playerC.PlayerExitHandler(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	expected := struct {
		message string
		Data    models.Player
	}{}
	json.Unmarshal(resp, &expected)
	assert.Equal(t, 200, w.Result().StatusCode)

}
