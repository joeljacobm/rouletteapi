package controllers

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rouletteapi/appconfigs"
	"rouletteapi/models"
	"rouletteapi/prometheus"
)

func writeJSON(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}

func writeErrorWithMsg(w http.ResponseWriter, err error, statuscode int) {
	// statusCode := http.statusCode
	prometheus.ErrorCounter.Inc()

	http.Error(w, err.Error(), statuscode)
}

func checkError(w http.ResponseWriter, err error) {
	if err == sql.ErrNoRows {
		writeErrorWithMsg(w, errors.New("Requested resource not found"), http.StatusNotFound)
		return
	}
	writeErrorWithMsg(w, err, http.StatusInternalServerError)
}

func generateHash() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	s := fmt.Sprintf("%X", b)

	return s
}

func calculateResult(bets []models.Bet) ([]models.Bet, error) {

	for k := range bets {
		switch bets[k].BetType {
		case 1:
			if bets[k].Selection == bets[k].BetResult.Number {
				oddsDecimal := appconfigs.GetRouletteOddsMap(bets[k].BetType)
				resultWin(&bets[k], oddsDecimal)
			} else {
				resultLost(&bets[k])
			}
		case 2:
			if bets[k].Selection == bets[k].BetResult.Colour {
				oddsDecimal := appconfigs.GetRouletteOddsMap(bets[k].BetType)
				resultWin(&bets[k], oddsDecimal)

			} else {
				resultLost(&bets[k])
			}
		case 3:
			var result int
			if bets[k].BetResult.Number%2 == 0 {
				result = 1
				bets[k].BetResult.OddEven = 1
			} else {
				result = 2
				bets[k].BetResult.OddEven = 2
			}
			if bets[k].Selection == result {
				oddsDecimal := appconfigs.GetRouletteOddsMap(bets[k].BetType)
				resultWin(&bets[k], oddsDecimal)
			} else {
				resultLost(&bets[k])
			}
		default:
			resultLost(&bets[k])

		}
	}
	return bets, nil
}

func resultWin(bet *models.Bet, odds float64) {
	bet.TotalReturn = (odds + 1) * bet.Stake
	bet.Result = 1
	bet.ResultText = "WIN"
	bet.Odds = odds
}

func resultLost(bet *models.Bet) {
	bet.Result = 0
	bet.ResultText = "LOST"
}

func calculateLiability(bettype int, stake float64) float64 {
	odds := appconfigs.GetRouletteOddsMap(bettype)
	return odds * stake
}
