package controllers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"rouletteapi/configs"
	"rouletteapi/models"
	"rouletteapi/prometheus"
)

func writeJSON(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		panic(err)
	}
}

func writeErrorWithMsg(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	prometheus.ErrorCounter.Inc()

	http.Error(w, err.Error(), statusCode)
}

func GenerateHash() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)

	return s
}

func calculateResult(bets []models.Bet) ([]models.Bet, error) {

	for k := range bets {
		switch bets[k].BetType {
		case 1:
			if bets[k].Selection == bets[k].BetResult.Number {
				oddsDecimal := configs.GetRouletteOddsMap(bets[k].BetType)
				resultWin(&bets[k], oddsDecimal)
			} else {
				resultLost(&bets[k])
			}
		case 2:
			if bets[k].Selection == bets[k].BetResult.Colour {
				oddsDecimal := configs.GetRouletteOddsMap(bets[k].BetType)
				resultWin(&bets[k], oddsDecimal)

			} else {
				resultLost(&bets[k])
			}
		case 3:
			var result int
			if bets[k].BetResult.Number%2 == 0 {
				result = 1
			} else {
				result = 2
			}
			if bets[k].Selection == result {
				oddsDecimal := configs.GetRouletteOddsMap(bets[k].BetType)
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
