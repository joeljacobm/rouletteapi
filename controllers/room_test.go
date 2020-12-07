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

func TestRoomHandlerDefaultVariant(t *testing.T) {
	testData := struct {
		VariantType int `json:"variant_type"`
	}{
		VariantType: 0,
	}

	appconfigs.LoadRouletteVariantMap("../appconfigs/variantconfig.json")

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
	assert.Equal(t, 7, expected.Data.RoomVariant.MaxPlayers)
	assert.Equal(t, 10, expected.Data.RoomVariant.NumOfRounds)
	assert.Equal(t, 1, expected.Data.RoomVariant.VariantType)
	assert.Equal(t, "sky_roulette", expected.Data.RoomVariant.VariantName)
}

func TestRoomHandlerCustomVariant(t *testing.T) {

	testData := struct {
		VariantType int `json:"variant_type"`
	}{
		VariantType: 6,
	}
	appconfigs.LoadRouletteVariantMap("../mocks/testconfig.json")

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
	assert.Equal(t, 15, expected.Data.RoomVariant.MaxPlayers)
	assert.Equal(t, 10, expected.Data.RoomVariant.NumOfRounds)
	assert.Equal(t, 6, expected.Data.RoomVariant.VariantType)
	assert.Equal(t, "deal_or_no_deal", expected.Data.RoomVariant.VariantName)
}

func TestRoomGetAllHandler(t *testing.T) {

	expectedData := []models.Room{{ID: "5F9E314910", RoomVariant: models.Variant{VariantType: 1, VariantName: "sky_roulette"}},
		{ID: "13D66E88E2", RoomVariant: models.Variant{VariantType: 2, VariantName: "superboost_roulett"}}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/room", nil)
	roomC.RoomGetAllHandler(w, r)
	resp, _ := ioutil.ReadAll(w.Body)

	actual := struct {
		message string
		Data    []models.Room
	}{}
	json.Unmarshal(resp, &actual)

	room1 := actual.Data[0]
	room2 := actual.Data[1]

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, expectedData[0].ID, room1.ID)
	assert.Equal(t, expectedData[0].CurrentRound, room1.CurrentRound)
	assert.Equal(t, expectedData[0].RoomVariant.VariantType, room1.RoomVariant.VariantType)
	assert.Equal(t, expectedData[0].RoomVariant.VariantName, room1.RoomVariant.VariantName)
	assert.Equal(t, expectedData[1].ID, room2.ID)
	assert.Equal(t, expectedData[1].CurrentRound, room2.CurrentRound)
	assert.Equal(t, expectedData[1].RoomVariant.VariantType, room2.RoomVariant.VariantType)
	assert.Equal(t, expectedData[1].RoomVariant.VariantName, room2.RoomVariant.VariantName)
}

func TestRoomGetHandler(t *testing.T) {

	expectedData := models.Room{ID: "5F9E314910", RoomVariant: models.Variant{VariantType: 1, VariantName: "sky_roulette"}}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/room/{id}", nil)
	vars := map[string]string{
		"id": "5F9E314910",
	}
	r = mux.SetURLVars(r, vars)

	roomC.RoomGetHandler(w, r)
	resp, _ := ioutil.ReadAll(w.Body)

	actual := struct {
		message string
		Data    models.Room
	}{}
	json.Unmarshal(resp, &actual)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, expectedData.ID, actual.Data.ID)
	assert.Equal(t, expectedData.RoomVariant.VariantType, actual.Data.RoomVariant.VariantType)
	assert.Equal(t, expectedData.RoomVariant.VariantName, actual.Data.RoomVariant.VariantName)

}
