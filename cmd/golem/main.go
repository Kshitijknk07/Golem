package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Golem/internal/api"
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

	metricStorage := storage.NewMemoryStorage()

	collector := collector.NewCollector(metricStorage)
	go collector.Start(ctx, 5*time.Second)

	apiServer := api.NewServer(metricStorage)
	server := &http.Server{
		Addr:    ":8080",
		Handler: apiServer.Router(),
	}

	// Check if static files directory exists
	staticDir := "web/static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Warning: Static files directory '%s' does not exist", staticDir)
	} else {
		log.Printf("Static files directory '%s' found", staticDir)
	}

	go func() {
		log.Println("Starting Golem monitoring server on http://localhost:8080")
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
