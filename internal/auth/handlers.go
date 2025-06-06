package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	UserStore  UserStorage
	JWTService *JWTService
}

// RegisterHandler handles user registration
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req UserCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := h.UserStore.CreateUser(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// LoginHandler handles user login and JWT issuance
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := h.UserStore.GetUserByUsername(req.Username)
	if err != nil || !user.IsActive {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if err := CheckPassword(req.Password, user.PasswordHash); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	h.UserStore.UpdateLastLogin(user.ID)
	token, expiresAt, err := h.JWTService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	resp := LoginResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: expiresAt,
	}
	json.NewEncoder(w).Encode(resp)
}

// ListUsersHandler returns all users (admin only)
func (h *Handler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserStore.ListUsers()
	if err != nil {
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// UpdateUserHandler updates a user (admin only)
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/auth/users/")
	var req UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := h.UserStore.UpdateUser(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// DeleteUserHandler deletes a user (admin only)
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/auth/users/")
	if err := h.UserStore.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CheckPassword compares a plaintext password with a hash
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
