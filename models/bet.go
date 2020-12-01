package models

type Bet struct {
	RoundNo     int     `json:"roundid"`
	BetType     int     `json:"bettype"`
	Stake       float64 `json:"stake"`
	Odds        float64 `json:"odds"`
	Liability   float64 `json:"liability"`
	Selection   int     `json:"selection"`
	Result      int     `json:"result"`
	TotalReturn float64 `json:"total_return"`
	BetOutcome  OutCome `json:"bet_outcome"`
}

type OutCome struct {
	Number  int `json:"number"`
	Colour  int `json:"colour"`
	OddEven int `json:"oddeven"`
}

type BetType struct {
	BetName     string  `json:"betname"`
	OddsDecimal float64 `json:"oddsdecimal"`
}
