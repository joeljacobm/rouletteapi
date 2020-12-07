package routes

import (
	"net/http"
	_ "net/http/pprof"
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
	r.PathPrefix("/debug/").Handler(http.DefaultServeMux)
	r.HandleFunc("/room", roomController.RoomHandler).Methods("POST")
	r.HandleFunc("/room", roomController.RoomGetAllHandler).Methods("GET")
	r.HandleFunc("/room/variants", roomController.RoomVariantHandler).Methods("GET")
	r.HandleFunc("/room/{id}", roomController.RoomGetHandler).Methods("GET")

	r.HandleFunc("/player", playerController.PlayerHandler).Methods("GET")
	r.HandleFunc("/player/result/{playerid}/{roomid}/{roundno}", playerController.PlayerGetResultHandler).Methods("GET")
	r.HandleFunc("/player/join", playerController.PlayerJoinHandler).Methods("POST")
	r.HandleFunc("/player/bet", playerController.PlayerBetHandler).Methods("POST")
	r.HandleFunc("/player/ready", playerController.PlayerReadyHandler).Methods("POST")
	r.HandleFunc("/player/exit", playerController.PlayerExitHandler).Methods("POST")

	r.HandleFunc("/bet/{roomid}", betController.BetHandler).Methods("GET")
	r.HandleFunc("/bet/result", betController.BetResultHandler).Methods("POST")

	return r
}
