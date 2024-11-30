package api

import (
	"encoding/json"
	"golem/internal/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Fetch Metrics
func FetchMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := models.GetAllMetrics()
	log.Println("Fetched Metrics:", metrics)
	json.NewEncoder(w).Encode(metrics)
}

// Save Metrics
func SaveMetrics(w http.ResponseWriter, r *http.Request) {
	var metric models.Metric
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		log.Println("Invalid JSON data:", err)
		return
	}
	log.Println("Received Metric:", metric)
	models.SaveMetric(metric)
	json.NewEncoder(w).Encode(metric)
}

// Update Metrics
func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	var metric models.Metric
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		log.Println("Invalid JSON data:", err)
		return
	}
	log.Println("Updated Metric:", metric)
	models.UpdateMetric(metric)
	json.NewEncoder(w).Encode(metric)
}

// Routes
func RegisterMetricsRoutes(router *mux.Router) {
	router.HandleFunc("/metrics", FetchMetrics).Methods("GET")
	router.HandleFunc("/metrics", SaveMetrics).Methods("POST")
	router.HandleFunc("/metrics", UpdateMetrics).Methods("PUT")
}
