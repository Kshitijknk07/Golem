package api

import (
	"encoding/json"
	"golem/internal/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func FetchMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := models.GetAllMetrics()
	if err != nil {
		log.Printf("Failed to fetch metrics: %v", err)
		http.Error(w, "Failed to fetch metrics", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		log.Printf("Failed to encode metrics: %v", err)
		http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
	}
}

func SaveMetrics(w http.ResponseWriter, r *http.Request) {
	var metric models.Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		log.Printf("Invalid JSON data: %v", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if err := models.SaveMetric(metric); err != nil {
		log.Printf("Failed to save metric: %v", err)
		http.Error(w, "Failed to save metric", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(metric); err != nil {
		log.Printf("Failed to encode metric: %v", err)
		http.Error(w, "Failed to encode metric", http.StatusInternalServerError)
	}
}

func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	var metric models.Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		log.Printf("Invalid JSON data: %v", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if err := models.UpdateMetric(metric); err != nil {
		log.Printf("Failed to update metric: %v", err)
		http.Error(w, "Failed to update metric", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(metric); err != nil {
		log.Printf("Failed to encode metric: %v", err)
		http.Error(w, "Failed to encode metric", http.StatusInternalServerError)
	}
}

func RegisterMetricsRoutes(router *mux.Router) {
	router.HandleFunc("/metrics", FetchMetrics).Methods("GET")
	router.HandleFunc("/metrics", SaveMetrics).Methods("POST")
	router.HandleFunc("/metrics", UpdateMetrics).Methods("PUT")
}
