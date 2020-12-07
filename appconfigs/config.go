package appconfigs

import (
	"encoding/json"
	"log"
	"os"
	"rouletteapi/models"
)

var (
	rouletteVariantMap map[int]models.Variant
	rouletteBetOddsMap map[int]models.BetType
)

// LoadRouletteVariantMap loads variantconfig.json file into the map
func LoadRouletteVariantMap(file string) {

	rouletteVariantMap = make(map[int]models.Variant)
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Cannot load roulette variant config with error %s", err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&rouletteVariantMap)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully loaded Roulette Variant Config")
}

func GetRouletteVariantMap(variantType int) models.Variant {

	if _, ok := rouletteVariantMap[variantType]; ok {
		return rouletteVariantMap[variantType]
	}

	return rouletteVariantMap[1] // Defaulting to the sky_roulette variant

}

func GetAllRouletteVariantMap() []models.Variant {

	var variants []models.Variant
	for _, variant := range rouletteVariantMap {
		variants = append(variants, variant)
	}

	return variants
}

// LoadRouletteOddsMap loads oddsconfig.json file into the map
func LoadRouletteOddsMap(file string) {
	rouletteBetOddsMap = make(map[int]models.BetType)
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Cannot load roulette odds config with error %s", err)

	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&rouletteBetOddsMap)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully loaded odds config")
}

func GetRouletteOddsMap(betType int) float64 {

	if _, ok := rouletteBetOddsMap[betType]; ok {
		return rouletteBetOddsMap[betType].OddsDecimal // Straight Up
	}
	return 1.00
}
