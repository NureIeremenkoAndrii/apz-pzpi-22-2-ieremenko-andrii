package models

import (
	"time"

	"github.com/google/uuid"
)

// Room represents a room in the system
type Room struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateRoomRequest represents a request to create a new room
type CreateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RoomListResponse represents a response for listing rooms
type RoomListResponse struct {
	Rooms []Room `json:"rooms"`
	Total int    `json:"total"`
}
