package api

import (
	"encoding/json"
	"net/http"
	"time"

	"Golem/internal/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	storage storage.MetricStorage
}

func NewServer(storage storage.MetricStorage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/metrics", s.getLatestMetrics).Methods("GET")
	r.HandleFunc("/api/metrics/history", s.getMetricsHistory).Methods("GET")

	fs := http.FileServer(http.Dir("web/static"))
	r.PathPrefix("/").Handler(fs)

	return r
}

func (s *Server) getLatestMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := s.storage.GetLatestMetrics()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (s *Server) getMetricsHistory(w http.ResponseWriter, r *http.Request) {
	duration := 1 * time.Hour

	durationParam := r.URL.Query().Get("duration")
	if durationParam != "" {
		parsedDuration, err := time.ParseDuration(durationParam)
		if err == nil {
			duration = parsedDuration
		}
	}

	metrics := s.storage.GetMetricsHistory(duration)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
