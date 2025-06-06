package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Golem/internal/auth"
	"Golem/internal/collector"
	"Golem/internal/metrics"
	"Golem/internal/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	storage              storage.MetricStorage
	healthCheckStorage   storage.HealthCheckStorage
	healthCheckCollector *collector.HealthCheckCollector

	userStorage auth.UserStorage
	jwtService  *auth.JWTService
	authHandler *auth.Handler
}

func NewServer(storage storage.MetricStorage, healthCheckStorage storage.HealthCheckStorage, healthCheckCollector *collector.HealthCheckCollector, userStorage auth.UserStorage, jwtService *auth.JWTService) *Server {
	return &Server{
		storage:              storage,
		healthCheckStorage:   healthCheckStorage,
		healthCheckCollector: healthCheckCollector,
		userStorage:          userStorage,
		jwtService:           jwtService,
		authHandler:          &auth.Handler{UserStore: userStorage, JWTService: jwtService},
	}
}

func (s *Server) Router() http.Handler {
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/api/auth/register", s.authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", s.authHandler.LoginHandler).Methods("POST")

	// User management (admin only)
	userSubrouter := r.PathPrefix("/api/auth/users").Subrouter()
	userSubrouter.Use(auth.JWTAuthMiddleware(s.jwtService))
	userSubrouter.Use(auth.RequireRoleMiddleware(auth.RoleAdmin))
	userSubrouter.HandleFunc("", s.authHandler.ListUsersHandler).Methods("GET")
	userSubrouter.HandleFunc("/{id}", s.authHandler.UpdateUserHandler).Methods("PUT")
	userSubrouter.HandleFunc("/{id}", s.authHandler.DeleteUserHandler).Methods("DELETE")

	r.HandleFunc("/api/metrics", s.getLatestMetrics).Methods("GET")
	r.HandleFunc("/api/metrics/history", s.getMetricsHistory).Methods("GET")

	r.HandleFunc("/api/health-checks", s.getHealthChecks).Methods("GET")
	r.HandleFunc("/api/health-checks", s.createHealthCheck).Methods("POST")
	r.HandleFunc("/api/health-checks/{id}", s.getHealthCheck).Methods("GET")
	r.HandleFunc("/api/health-checks/{id}", s.updateHealthCheck).Methods("PUT")
	r.HandleFunc("/api/health-checks/{id}", s.deleteHealthCheck).Methods("DELETE")
	r.HandleFunc("/api/health-checks/{id}/history", s.getHealthCheckHistory).Methods("GET")

	fs := http.FileServer(http.Dir("web/static"))
	r.PathPrefix("/").Handler(fs)

	return r
}

func (s *Server) getLatestMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := s.storage.GetLatestMetrics()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get latest metrics: %v", err), http.StatusInternalServerError)
		return
	}

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

	metrics, err := s.storage.GetMetricsHistory(duration)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get metrics history: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (s *Server) getHealthChecks(w http.ResponseWriter, r *http.Request) {
	results, err := s.healthCheckStorage.GetAllHealthCheckResults()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get health checks: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *Server) getHealthCheck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := s.healthCheckStorage.GetHealthCheckResult(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) createHealthCheck(w http.ResponseWriter, r *http.Request) {
	var config metrics.HealthCheckConfig

	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if config.Interval == 0 {
		config.Interval = 60 * time.Second
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	if config.ID == "" {
		config.ID = fmt.Sprintf("check_%d", time.Now().UnixNano())
	}

	config.Enabled = true
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	err = s.healthCheckCollector.AddHealthCheck(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

func (s *Server) updateHealthCheck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var config metrics.HealthCheckConfig

	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.ID = id
	config.UpdatedAt = time.Now()

	err = s.healthCheckCollector.UpdateHealthCheck(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (s *Server) deleteHealthCheck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := s.healthCheckCollector.DeleteHealthCheck(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getHealthCheckHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	duration := 24 * time.Hour

	durationParam := r.URL.Query().Get("duration")
	if durationParam != "" {
		parsedDuration, err := time.ParseDuration(durationParam)
		if err == nil {
			duration = parsedDuration
		}
	}

	history, err := s.healthCheckStorage.GetHealthCheckHistory(id, duration)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get health check history: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
