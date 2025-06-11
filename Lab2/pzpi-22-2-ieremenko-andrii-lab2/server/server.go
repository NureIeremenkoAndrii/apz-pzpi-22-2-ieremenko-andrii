package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andrii/apz-pzpi-22-2-ieremenko-andrii/Lab2/pzpi-22-2-ieremenko-andrii-lab2/auth"
	_ "github.com/andrii/apz-pzpi-22-2-ieremenko-andrii/Lab2/pzpi-22-2-ieremenko-andrii-lab2/docs"
	"github.com/andrii/apz-pzpi-22-2-ieremenko-andrii/Lab2/pzpi-22-2-ieremenko-andrii-lab2/models"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server represents the HTTP server
type Server struct {
	users    map[uuid.UUID]*models.User
	roles    map[string]*models.Role
	rooms    map[uuid.UUID]*models.Room
	metrics  map[uuid.UUID]*models.Metric
	readings map[uuid.UUID][]*models.MetricReading
}

// NewServer creates a new server instance
func NewServer() *Server {
	// Initialize default roles
	roles := map[string]*models.Role{
		"admin": {
			ID:          uuid.New(),
			Name:        "admin",
			Description: "Administrator role with full access",
			Permissions: []string{"read", "write", "delete", "manage_users", "manage_roles", "manage_metrics", "manage_rooms"},
		},
		"user": {
			ID:          uuid.New(),
			Name:        "user",
			Description: "Regular user role",
			Permissions: []string{"read", "write"},
		},
	}

	return &Server{
		users:    make(map[uuid.UUID]*models.User),
		roles:    roles,
		rooms:    make(map[uuid.UUID]*models.Room),
		metrics:  make(map[uuid.UUID]*models.Metric),
		readings: make(map[uuid.UUID][]*models.MetricReading),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /register [post]
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if username is already taken
	for _, user := range s.users {
		if user.Username == req.Username {
			http.Error(w, "Username already taken", http.StatusBadRequest)
			return
		}
	}

	// Create user with default role if none specified
	userRoles := make([]models.Role, 0)
	if len(req.Roles) == 0 {
		userRoles = append(userRoles, *s.roles["user"])
	} else {
		for _, roleName := range req.Roles {
			if role, exists := s.roles[roleName]; exists {
				userRoles = append(userRoles, *role)
			}
		}
	}

	now := time.Now()
	user := &models.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Password:  req.Password, // In a real application, this should be hashed
		Email:     req.Email,
		Roles:     userRoles,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[user.ID] = user

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Login godoc
// @Summary Login user
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /login [post]
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by username
	var user *models.User
	for _, u := range s.users {
		if u.Username == req.Username {
			user = u
			break
		}
	}

	if user == nil || user.Password != req.Password { // In a real application, compare hashed passwords
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of all registered users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.UserListResponse
// @Router /users [get]
func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Check if user has admin role
	claims, err := auth.ValidateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Find user by username
	var user *models.User
	for _, u := range s.users {
		if u.Username == claims.Username {
			user = u
			break
		}
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has admin role
	hasAdminRole := false
	for _, role := range user.Roles {
		if role.Name == "admin" {
			hasAdminRole = true
			break
		}
	}

	if !hasAdminRole {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	users := make([]models.User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, *u)
	}

	json.NewEncoder(w).Encode(models.UserListResponse{
		Users: users,
		Total: len(users),
	})
}

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role with specified permissions
// @Tags roles
// @Accept json
// @Produce json
// @Param request body models.CreateRoleRequest true "Role creation request"
// @Success 200 {object} models.Role
// @Failure 400 {object} map[string]string
// @Router /roles [post]
func (s *Server) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Check if user has admin role
	claims, err := auth.ValidateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Find user by username
	var user *models.User
	for _, u := range s.users {
		if u.Username == claims.Username {
			user = u
			break
		}
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has admin role
	hasAdminRole := false
	for _, role := range user.Roles {
		if role.Name == "admin" {
			hasAdminRole = true
			break
		}
	}

	if !hasAdminRole {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	var req models.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if role already exists
	if _, exists := s.roles[req.Name]; exists {
		http.Error(w, "Role already exists", http.StatusBadRequest)
		return
	}

	role := &models.Role{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
	}

	s.roles[role.Name] = role

	json.NewEncoder(w).Encode(role)
}

// GetUserIDFromToken extracts user ID from JWT token
func (s *Server) GetUserIDFromToken(r *http.Request) (uuid.UUID, error) {
	claims, err := auth.ValidateToken(r)
	if err != nil {
		return uuid.Nil, err
	}

	// Find user by username
	for _, user := range s.users {
		if user.Username == claims.Username {
			return user.ID, nil
		}
	}

	return uuid.Nil, fmt.Errorf("user not found")
}

// CreateRoom godoc
// @Summary Create a new room
// @Description Create a new room with specified name and description
// @Tags rooms
// @Accept json
// @Produce json
// @Param request body models.CreateRoomRequest true "Room creation request"
// @Success 200 {object} models.Room
// @Failure 400 {object} map[string]string
// @Router /rooms [post]
func (s *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	// Check if user has permission to manage rooms
	userID, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, exists := s.users[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	hasPermission := false
	for _, role := range user.Roles {
		if role.Name == "admin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	var req models.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	room := &models.Room{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.rooms[room.ID] = room

	json.NewEncoder(w).Encode(room)
}

// ListRooms godoc
// @Summary List all rooms
// @Description Get a list of all rooms
// @Tags rooms
// @Accept json
// @Produce json
// @Success 200 {object} models.RoomListResponse
// @Router /rooms [get]
func (s *Server) ListRooms(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	_, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rooms := make([]models.Room, 0, len(s.rooms))
	for _, r := range s.rooms {
		rooms = append(rooms, *r)
	}

	json.NewEncoder(w).Encode(models.RoomListResponse{
		Rooms: rooms,
		Total: len(rooms),
	})
}

// GetRoom godoc
// @Summary Get room details
// @Description Get details of a specific room with its metrics
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} models.Room
// @Failure 404 {object} map[string]string
// @Router /rooms/{id} [get]
func (s *Server) GetRoom(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	_, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Path[len("/rooms/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	room, exists := s.rooms[id]
	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(room)
}

// DeleteRoom godoc
// @Summary Delete a room
// @Description Delete a room and all its metrics
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /rooms/{id} [delete]
func (s *Server) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	// Check if user has permission to manage rooms
	userID, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, exists := s.users[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	hasPermission := false
	for _, role := range user.Roles {
		if role.Name == "admin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	idStr := r.URL.Path[len("/rooms/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	if _, exists := s.rooms[id]; !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Delete all metrics associated with this room
	for metricID, metric := range s.metrics {
		if metric.RoomID == id {
			delete(s.metrics, metricID)
			delete(s.readings, metricID)
		}
	}

	delete(s.rooms, id)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Room deleted successfully",
	})
}

// CreateMetric godoc
// @Summary Create a new metric
// @Description Create a new household metric
// @Tags metrics
// @Accept json
// @Produce json
// @Param request body models.CreateMetricRequest true "Metric creation request"
// @Success 200 {object} models.Metric
// @Failure 400 {object} map[string]string
// @Router /metrics [post]
func (s *Server) CreateMetric(w http.ResponseWriter, r *http.Request) {
	// Check if user has permission to manage metrics
	userID, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, exists := s.users[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	hasPermission := false
	for _, role := range user.Roles {
		if role.Name == "admin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	var req models.CreateMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if room exists
	if _, exists := s.rooms[req.RoomID]; !exists {
		http.Error(w, "Room not found", http.StatusBadRequest)
		return
	}

	now := time.Now()
	metric := &models.Metric{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
		RoomID:      req.RoomID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.metrics[metric.ID] = metric
	s.readings[metric.ID] = make([]*models.MetricReading, 0)

	json.NewEncoder(w).Encode(metric)
}

// ListMetrics godoc
// @Summary List all metrics
// @Description Get a list of all metrics
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} models.MetricListResponse
// @Router /metrics [get]
func (s *Server) ListMetrics(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	_, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	metrics := make([]models.Metric, 0, len(s.metrics))
	for _, m := range s.metrics {
		metrics = append(metrics, *m)
	}

	json.NewEncoder(w).Encode(models.MetricListResponse{
		Metrics: metrics,
		Total:   len(metrics),
	})
}

// GetMetric godoc
// @Summary Get metric details
// @Description Get details of a specific metric with its readings
// @Tags metrics
// @Accept json
// @Produce json
// @Param id path string true "Metric ID"
// @Success 200 {object} models.MetricWithReadings
// @Failure 404 {object} map[string]string
// @Router /metrics/{id} [get]
func (s *Server) GetMetric(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	_, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Path[len("/metrics/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid metric ID", http.StatusBadRequest)
		return
	}

	metric, exists := s.metrics[id]
	if !exists {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	readings := make([]models.MetricReading, 0, len(s.readings[id]))
	for _, r := range s.readings[id] {
		readings = append(readings, *r)
	}

	json.NewEncoder(w).Encode(models.MetricWithReadings{
		Metric:   *metric,
		Readings: readings,
	})
}

// DeleteMetric godoc
// @Summary Delete a metric
// @Description Delete a metric and all its readings
// @Tags metrics
// @Accept json
// @Produce json
// @Param id path string true "Metric ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /metrics/{id} [delete]
func (s *Server) DeleteMetric(w http.ResponseWriter, r *http.Request) {
	// Check if user has permission to manage metrics
	userID, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, exists := s.users[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	hasPermission := false
	for _, role := range user.Roles {
		if role.Name == "admin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	idStr := r.URL.Path[len("/metrics/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid metric ID", http.StatusBadRequest)
		return
	}

	if _, exists := s.metrics[id]; !exists {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	delete(s.metrics, id)
	delete(s.readings, id)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Metric deleted successfully",
	})
}

// AddReading godoc
// @Summary Add a reading
// @Description Add a new reading for a metric
// @Tags metrics
// @Accept json
// @Produce json
// @Param id path string true "Metric ID"
// @Param request body models.AddReadingRequest true "Reading request"
// @Success 200 {object} models.MetricReading
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /metrics/{id}/readings [post]
func (s *Server) AddReading(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	_, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Path[len("/metrics/"):]
	idStr = idStr[:len(idStr)-len("/readings")]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid metric ID", http.StatusBadRequest)
		return
	}

	if _, exists := s.metrics[id]; !exists {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	var req models.AddReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// If timestamp is not provided, use current time
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	reading := &models.MetricReading{
		ID:        uuid.New(),
		MetricID:  id,
		Value:     req.Value,
		Timestamp: req.Timestamp,
		CreatedAt: time.Now(),
	}

	s.readings[id] = append(s.readings[id], reading)

	json.NewEncoder(w).Encode(reading)
}

// GetReadings godoc
// @Summary Get metric readings
// @Description Get all readings for a metric
// @Tags metrics
// @Accept json
// @Produce json
// @Param id path string true "Metric ID"
// @Success 200 {object} models.ReadingListResponse
// @Failure 404 {object} map[string]string
// @Router /metrics/{id}/readings [get]
func (s *Server) GetReadings(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	_, err := s.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Path[len("/metrics/"):]
	idStr = idStr[:len(idStr)-len("/readings")]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid metric ID", http.StatusBadRequest)
		return
	}

	if _, exists := s.metrics[id]; !exists {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	readings := make([]models.MetricReading, 0, len(s.readings[id]))
	for _, r := range s.readings[id] {
		readings = append(readings, *r)
	}

	json.NewEncoder(w).Encode(models.ReadingListResponse{
		Readings: readings,
		Total:    len(readings),
	})
}

// Start starts the server
func (s *Server) Start(addr string) error {
	// API endpoints
	http.HandleFunc("/register", s.Register)
	http.HandleFunc("/login", s.Login)
	http.HandleFunc("/users", s.ListUsers)
	http.HandleFunc("/roles", s.CreateRole)

	// Room endpoints
	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateRoom(w, r)
		case http.MethodGet:
			s.ListRooms(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/rooms/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rooms/" {
			http.Error(w, "Invalid room ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			s.GetRoom(w, r)
		case http.MethodDelete:
			s.DeleteRoom(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Metric endpoints
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateMetric(w, r)
		case http.MethodGet:
			s.ListMetrics(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/metrics/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics/" {
			http.Error(w, "Invalid metric ID", http.StatusBadRequest)
			return
		}

		if r.URL.Path[len(r.URL.Path)-len("/readings"):] == "/readings" {
			switch r.Method {
			case http.MethodPost:
				s.AddReading(w, r)
			case http.MethodGet:
				s.GetReadings(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		switch r.Method {
		case http.MethodGet:
			s.GetMetric(w, r)
		case http.MethodDelete:
			s.DeleteMetric(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Swagger documentation
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return http.ListenAndServe(addr, nil)
}
