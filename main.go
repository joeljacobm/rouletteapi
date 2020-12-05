package main

import (
	"log"
	"net/http"
	"rouletteapi/configs"
	"rouletteapi/models"
	"rouletteapi/postgres"
	"rouletteapi/routes"
)

func main() {

	log.Println("Roulette API Server")
	configs.LoadRouletteVariantMap("configs/config.json")
	configs.LoadRouletteOddsMap("configs/oddsconfig.json")

	db := postgres.DefaultConnection()

	services := models.NewServices(db)

	router := routes.SetRoutes(services)

	log.Fatal(http.ListenAndServe(":8080", router))
}
