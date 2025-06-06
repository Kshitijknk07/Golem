package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
)

// UserStorage defines the interface for user storage operations
type UserStorage interface {
	CreateUser(user *UserCreate) (*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(id string, update *UserUpdate) (*User, error)
	DeleteUser(id string) error
	ListUsers() ([]*User, error)
	UpdateLastLogin(id string) error
}

// SQLiteUserStorage implements UserStorage using SQLite
type SQLiteUserStorage struct {
	db *sql.DB
}

// NewSQLiteUserStorage creates a new SQLiteUserStorage instance
func NewSQLiteUserStorage(db *sql.DB) (*SQLiteUserStorage, error) {
	// Create users table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			last_login TIMESTAMP,
			is_active BOOLEAN NOT NULL DEFAULT true
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteUserStorage{db: db}, nil
}

// CreateUser creates a new user
func (s *SQLiteUserStorage) CreateUser(user *UserCreate) (*User, error) {
	// Check if username already exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", user.Username).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	newUser := &User{
		ID:           uuid.New().String(),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: string(hashedPassword),
		Role:         user.Role,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}

	_, err = s.db.Exec(`
		INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, newUser.ID, newUser.Username, newUser.Email, newUser.PasswordHash, newUser.Role,
		newUser.CreatedAt, newUser.UpdatedAt, newUser.IsActive)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// GetUserByID retrieves a user by ID
func (s *SQLiteUserStorage) GetUserByID(id string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(`
		SELECT id, username, email, password_hash, role, created_at, updated_at, last_login, is_active
		FROM users WHERE id = ?
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *SQLiteUserStorage) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(`
		SELECT id, username, email, password_hash, role, created_at, updated_at, last_login, is_active
		FROM users WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user's information
func (s *SQLiteUserStorage) UpdateUser(id string, update *UserUpdate) (*User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if update.Email != nil {
		user.Email = *update.Email
	}
	if update.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*update.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
	}
	if update.Role != nil {
		user.Role = *update.Role
	}
	if update.IsActive != nil {
		user.IsActive = *update.IsActive
	}

	user.UpdatedAt = time.Now()

	_, err = s.db.Exec(`
		UPDATE users
		SET email = ?, password_hash = ?, role = ?, updated_at = ?, is_active = ?
		WHERE id = ?
	`, user.Email, user.PasswordHash, user.Role, user.UpdatedAt, user.IsActive, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *SQLiteUserStorage) DeleteUser(id string) error {
	result, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// ListUsers retrieves all users
func (s *SQLiteUserStorage) ListUsers() ([]*User, error) {
	rows, err := s.db.Query(`
		SELECT id, username, email, password_hash, role, created_at, updated_at, last_login, is_active
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateLastLogin updates the last login timestamp for a user
func (s *SQLiteUserStorage) UpdateLastLogin(id string) error {
	_, err := s.db.Exec("UPDATE users SET last_login = ? WHERE id = ?", time.Now(), id)
	return err
}
