package api

import (
	"golem/internal/models"
	"net/http"

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		metrics := models.GetAllMetrics()
		if err := conn.WriteJSON(metrics); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func RegisterWebSocketRoutes(router *mux.Router) {
	router.HandleFunc("/ws/metrics", StreamMetrics).Methods("GET")
}
