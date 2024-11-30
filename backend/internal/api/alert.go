package api

import (
	"encoding/json"
	"golem/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

// Trigger Alerts
func TriggerAlert(w http.ResponseWriter, r *http.Request) {
	var alert models.Alert
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	models.TriggerAlert(alert)
	json.NewEncoder(w).Encode(alert)
}

// Manage Alerts
func ManageAlert(w http.ResponseWriter, r *http.Request) {
	alerts := models.GetAllAlerts()
	json.NewEncoder(w).Encode(alerts)
}

// Routes
func RegisterAlertRoutes(router *mux.Router) {
	router.HandleFunc("/alerts", TriggerAlert).Methods("POST")
	router.HandleFunc("/alerts", ManageAlert).Methods("GET")
}
