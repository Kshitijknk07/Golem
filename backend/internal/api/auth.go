package api

import (
	"encoding/json"
	"golem/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Dummy check - replace with real authentication
	if creds.Username != "user" || creds.Password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// RegisterAuthRoutes registers authentication routes
func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/login", LoginHandler).Methods("POST")
}
