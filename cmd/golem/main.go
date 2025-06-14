package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"Golem/internal/api"
	"Golem/internal/auth"
	"Golem/internal/collector"
	"Golem/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: Could not get current working directory: %v", err)
	} else {
		log.Printf("Current working directory: %s", cwd)
	}

	// Create data directory if it doesn't exist
	dataDir := filepath.Join(cwd, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize SQLite storage for metrics and health checks
	dbPath := filepath.Join(dataDir, "golem.db")
	metricStorage, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite storage: %v", err)
	}
	defer metricStorage.Close()

	// Initialize SQLite DB for users (reuse golem.db)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open SQLite DB for users: %v", err)
	}
	defer db.Close()

	userStorage, err := auth.NewSQLiteUserStorage(db)
	if err != nil {
		log.Fatalf("Failed to initialize user storage: %v", err)
	}

	// JWT secret and duration (should be from env in production)
	jwtSecret := "supersecretkey" // TODO: use os.Getenv in production
	tokenDuration := 24 * time.Hour
	jwtService := auth.NewJWTService(jwtSecret, tokenDuration)

	collector := collector.NewCollector(metricStorage)
	go collector.Start(ctx, 5*time.Second)

	healthCheckCollector := collector.NewHealthCheckCollector(metricStorage)
	go healthCheckCollector.Start(ctx)

	apiServer := api.NewServer(metricStorage, metricStorage, healthCheckCollector, userStorage, jwtService)
	server := &http.Server{
		Addr:    ":8899",
		Handler: apiServer.Router(),
	}

	staticDir := "web/static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Warning: Static files directory '%s' does not exist", staticDir)
	} else {
		log.Printf("Static files directory '%s' found", staticDir)
	}

	go func() {
		log.Println("Starting Golem monitoring server on http://localhost:8899")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
