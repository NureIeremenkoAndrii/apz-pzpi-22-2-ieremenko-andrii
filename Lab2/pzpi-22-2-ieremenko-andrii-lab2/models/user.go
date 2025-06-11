package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
// @Description User information
type User struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username  string    `json:"username" example:"johndoe"`
	Password  string    `json:"-"` // Password is not exposed in JSON
	Email     string    `json:"email" example:"john@example.com"`
	Roles     []Role    `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginRequest represents the login request payload
// @Description Login request payload
type LoginRequest struct {
	Username string `json:"username" example:"johndoe" binding:"required"`
	Password string `json:"password" example:"secretpassword" binding:"required"`
}

// RegisterRequest represents the registration request payload
// @Description Registration request payload
type RegisterRequest struct {
	Username string   `json:"username" example:"johndoe" binding:"required"`
	Password string   `json:"password" example:"secretpassword" binding:"required"`
	Email    string   `json:"email" example:"john@example.com" binding:"required,email"`
	Roles    []string `json:"roles" example:"['user', 'admin']"`
}

// AuthResponse represents the authentication response
// @Description Authentication response containing JWT token and user information
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// UserListResponse represents the response for listing users
// @Description Response containing a list of users
type UserListResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}
