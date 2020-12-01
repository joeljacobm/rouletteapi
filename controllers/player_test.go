package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rouletteapi/mocks"
	"rouletteapi/models"
	"rouletteapi/postgres"
	"testing"

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

	db := postgres.DefaultConnection()

	services := mocks.NewMockServices(db)
	playerC := NewPlayer(services.Player, services.Room, services.Bet, nil)
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

	assert.Equal(t, true, expected.Data.InRoom)
	assert.Equal(t, false, expected.Data.ReadyStatus)
	assert.Equal(t, "FA399BB381", expected.Data.RoomID)
	assert.Equal(t, "5EB8E63F56", expected.Data.ID)
}
