package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rouletteapi/models"
)

var (
	rouletteVariantMap map[int]models.Room
	rouletteBetOddsMap map[int]models.BetType
)

func LoadRouletteVariantMap(file string) {

	rouletteVariantMap = make(map[int]models.Room)
	f, err := os.Open(file)
	if err != nil {
		log.Println("Using the default config...")
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&rouletteVariantMap)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully loaded .config")
}

func GetRouletteVariantMap(variantType int) models.Room {

	if _, ok := rouletteVariantMap[variantType]; ok {
		return rouletteVariantMap[variantType]
	}

	return rouletteVariantMap[1] // Defaulting to the sky_roulette variant

}

func LoadRouletteOddsMap(file string) {
	rouletteBetOddsMap = make(map[int]models.BetType)
	f, err := os.Open(file)
	if err != nil {
		log.Println("Using the default config...")
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&rouletteBetOddsMap)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully loaded .config")
	fmt.Println(rouletteBetOddsMap)
}

func GetRouletteOddsMap(betType int) float64 {

	if _, ok := rouletteBetOddsMap[betType]; ok {

		if betType <= 36 {
			return rouletteBetOddsMap[1].OddsDecimal // Straight Up
		}
		return rouletteBetOddsMap[betType].OddsDecimal // Other Bet Types
	}

	return 1.00
}
