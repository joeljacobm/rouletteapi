package routes

import (
	"net/http"
	"rouletteapi/controllers"
	"rouletteapi/models"
	"rouletteapi/prometheus"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetRoutes(services models.Services) http.Handler {

	r := mux.NewRouter()

	roomController := controllers.NewRoom(services.Room)
	playerController := controllers.NewPlayer(services.Player, services.Room, services.Bet)
	betController := controllers.NewBet(services.Player, services.Room, services.Bet)

	r.Use(prometheus.PrometheusMiddleware)
	r.Handle("/metrics", promhttp.Handler())

	r.HandleFunc("/room", roomController.RoomHandler).Methods("POST")

	r.HandleFunc("/player/result/{playerid}/{roomid}/{roundno}", playerController.PlayerGetBetHandler).Methods("GET")
	r.HandleFunc("/player/join", playerController.PlayerJoinHandler).Methods("POST")
	r.HandleFunc("/player/bet", playerController.PlayerBetHandler).Methods("POST")
	r.HandleFunc("/player/ready", playerController.PlayerReadyHandler).Methods("POST")

	r.HandleFunc("/bet/result", betController.BetResultHandler).Methods("POST")

	return r
}
