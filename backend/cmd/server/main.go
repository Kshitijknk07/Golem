package main

import (
	"golem/internal/api"
	"golem/internal/db"
	service "golem/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	db.InitDB()
	service.InitMetricsCollection()

	router := mux.NewRouter()

	api.RegisterMetricsRoutes(router)
	api.RegisterAlertRoutes(router)
	api.RegisterWebSocketRoutes(router)
	api.RegisterAuthRoutes(router)
	router.Handle("/metrics", promhttp.Handler())

	log.Println("Server is running on port 4000")
	log.Fatal(http.ListenAndServe(":4000", router))
}
