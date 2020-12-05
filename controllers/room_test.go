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

	"github.com/stretchr/testify/assert"
)

func TestRoomHandlerDefaultVariant(t *testing.T) {
	testData := struct {
		VariantType int `json:"variant_type"`
	}{
		VariantType: 0,
	}

	configs.LoadRouletteVariantMap("../configs/config.json")

	services := mocks.NewMockServices(db)
	roomC := NewRoom(services.Room)

	b, _ := json.Marshal(testData)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/room", bytes.NewReader(b))
	r.Header.Add("Content-Type", "application/json")
	roomC.RoomHandler(w, r)
	resp, _ := ioutil.ReadAll(w.Body)

	expected := struct {
		message string
		Data    models.Room
	}{}
	json.Unmarshal(resp, &expected)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, 10, len(expected.Data.ID))
	assert.Equal(t, 7, expected.Data.MaxPlayers)
	assert.Equal(t, 10, expected.Data.NumOfRounds)
	assert.Equal(t, 1, expected.Data.VariantType)
	assert.Equal(t, "sky_roulette", expected.Data.VariantName)
}

func TestRoomHandlerCustomVariant(t *testing.T) {

	testData := struct {
		VariantType int `json:"variant_type"`
	}{
		VariantType: 6,
	}
	configs.LoadRouletteVariantMap("../mocks/testconfig.json")

	services := mocks.NewMockServices(db)
	roomC := NewRoom(services.Room)
	b, _ := json.Marshal(testData)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/room", bytes.NewReader(b))
	r.Header.Add("Content-Type", "application/json")
	roomC.RoomHandler(w, r)
	resp, _ := ioutil.ReadAll(w.Body)

	expected := struct {
		message string
		Data    models.Room
	}{}
	json.Unmarshal(resp, &expected)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, 10, len(expected.Data.ID))
	assert.Equal(t, 15, expected.Data.MaxPlayers)
	assert.Equal(t, 10, expected.Data.NumOfRounds)
	assert.Equal(t, 6, expected.Data.VariantType)
	assert.Equal(t, "deal_or_no_deal", expected.Data.VariantName)
}
