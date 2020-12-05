package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rouletteapi/mocks"
	"rouletteapi/models"
	"rouletteapi/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	db = postgres.DefaultConnection()
)

func TestBetResultHandler(t *testing.T) {
	testData := models.Bet{
		RoomID:  "100D2BF54B",
		RoundNo: 1,
	}

	services := mocks.NewMockServices(db)
	betC := NewBet(services.Player, services.Room, services.Bet)

	b, _ := json.Marshal(testData)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/bet/result", bytes.NewReader(b))
	betC.BetResultHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

}
