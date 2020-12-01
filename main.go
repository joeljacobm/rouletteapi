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
	playerController := controllers.NewPlayer(services.Player, services.Room, r)

	r.Use(prometheus.PrometheusMiddleware)
	r.HandleFunc("/room", roomController.RoomHandler).Methods("POST")
	
	r.HandleFunc("/player/join", playerController.PlayerJoinHandler).Methods("POST")

	r.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8080", r)
}
