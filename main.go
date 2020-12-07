package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rouletteapi/appconfigs"
	"rouletteapi/models"
	"rouletteapi/postgres"
	"rouletteapi/routes"
	"strconv"
	"syscall"
)

var (
	port int
)

func main() {

	flag.IntVar(&port, "port", 8080, "The port for running the HTTP server")
	flag.Parse()

	log.Println("-------Roulette API Server-------")
	appconfigs.LoadRouletteVariantMap("appconfigs/variantconfig.json")
	appconfigs.LoadRouletteOddsMap("appconfigs/oddsconfig.json")

	db := postgres.DefaultConnection()
	log.Println("Succefully connected to the database")
	services := models.NewServices(db)

	router := routes.SetRoutes(services)

	portStr := strconv.Itoa(port)
	go func() {
		err := http.ListenAndServe(":"+portStr, router)
		if err != nil {
			log.Fatalf("HTTP Server stopped with error %s", err)
		}
	}()

	log.Printf("Started HTTP server on port %d\n", port)
	exitCode := waitForStop()
	if exitCode != 0 {
		os.Exit(exitCode)
	}

}

func waitForStop() int {

	okSig := make(chan os.Signal, 2)
	signal.Notify(okSig, syscall.SIGTERM)

	failSig := make(chan os.Signal, 2)
	signal.Notify(failSig, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	for {
		select {
		case sig := <-okSig:
			log.Printf("Exit OK on Signal  %s", sig)
			return 0
		case sig := <-failSig:
			log.Printf("Exit FAIL on Signal  %s", sig)
			return 2
		}
	}

}
