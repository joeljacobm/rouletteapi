package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rouletteapi/models"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)


func TestBetResultHandler(t *testing.T) {
	testData := models.Bet{
		RoomID:  "100D2BF54B",
		RoundNo: 1,
	}


	b, _ := json.Marshal(testData)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/bet/result", bytes.NewReader(b))
	betC.BetResultHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

}

func TestBetHandler(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/bet/", nil)

	vars := map[string]string{
		"roomid": "100D2BF54B",
	}

	r = mux.SetURLVars(r, vars)
	r.Header.Add("Content-Type", "application/json")
	betC.BetHandler(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	expected := struct {
		message string
		Data    []models.Bet
	}{}
	json.Unmarshal(resp, &expected)

	bet1 := expected.Data[0]
	bet2 := expected.Data[1]
	json.Unmarshal(resp, &expected)
	assert.Equal(t, 1, bet1.BetType)
	assert.Equal(t, "100D2BF54B", bet1.RoomID)
	assert.Equal(t, 24, bet1.Selection)

	assert.Equal(t, 2, bet2.BetType)
	assert.Equal(t, 1, bet2.Selection)
	assert.Equal(t, "100D2BF54B", bet2.RoomID)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
