package api

import (
	"golem/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StreamMetrics(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		http.Error(w, "Failed to establish WebSocket connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics, err := models.GetAllMetrics()
			if err != nil {
				log.Printf("Failed to fetch metrics for WebSocket: %v", err)
				continue
			}

			if err := conn.WriteJSON(metrics); err != nil {
				log.Printf("Failed to write metrics to WebSocket: %v", err)
				return
			}
		}
	}
}

func RegisterWebSocketRoutes(router *mux.Router) {
	router.HandleFunc("/ws/metrics", StreamMetrics).Methods("GET")
}
