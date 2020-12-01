package main

import (
	"log"
	"net/http"
	"rouletteapi/configs"
	"rouletteapi/controllers"
	"rouletteapi/models"
	"rouletteapi/postgres"
	"rouletteapi/prometheus"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	log.Println("Roulette API Server")
	configs.LoadRouletteVariantMap("configs/config.json")
	configs.LoadRouletteOddsMap("configs/oddsconfig.json")
	router := mux.NewRouter()
	db := postgres.DefaultConnection()

	services := models.NewServices(db)

	r := router.PathPrefix("/api/v1").Subrouter()

	roomController := controllers.NewRoom(services.Room, r)
	playerController := controllers.NewPlayer(services.Player, services.Room, services.Bet, r)
	betController := controllers.NewBet(services.Player, services.Room, services.Bet, r)

	r.Use(prometheus.PrometheusMiddleware)
	r.Handle("/metrics", promhttp.Handler())

	r.HandleFunc("/room", roomController.RoomHandler).Methods("POST")

	r.HandleFunc("/player/bet", playerController.PlayerGetBetHandler).Methods("GET")
	r.HandleFunc("/player/join", playerController.PlayerJoinHandler).Methods("POST")
	r.HandleFunc("/player/bet", playerController.PlayerBetHandler).Methods("POST")
	r.HandleFunc("/player/ready", playerController.PlayerReadyHandler).Methods("POST")

	r.HandleFunc("/bet/result", betController.BetResultHandler).Methods("POST")

	http.ListenAndServe(":8080", router)
}
