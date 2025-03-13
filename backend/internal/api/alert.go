package api

import (
	"encoding/json"
	"golem/internal/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func TriggerAlert(w http.ResponseWriter, r *http.Request) {
	var alert models.Alert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		log.Printf("Invalid JSON data: %v", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if err := models.TriggerAlert(alert); err != nil {
		log.Printf("Failed to trigger alert: %v", err)
		http.Error(w, "Failed to trigger alert", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(alert); err != nil {
		log.Printf("Failed to encode alert: %v", err)
		http.Error(w, "Failed to encode alert", http.StatusInternalServerError)
	}
}

func ManageAlert(w http.ResponseWriter, r *http.Request) {
	alerts, err := models.GetAllAlerts()
	if err != nil {
		log.Printf("Failed to fetch alerts: %v", err)
		http.Error(w, "Failed to fetch alerts", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(alerts); err != nil {
		log.Printf("Failed to encode alerts: %v", err)
		http.Error(w, "Failed to encode alerts", http.StatusInternalServerError)
	}
}

func RegisterAlertRoutes(router *mux.Router) {
	router.HandleFunc("/alerts", TriggerAlert).Methods("POST")
	router.HandleFunc("/alerts", ManageAlert).Methods("GET")
}
