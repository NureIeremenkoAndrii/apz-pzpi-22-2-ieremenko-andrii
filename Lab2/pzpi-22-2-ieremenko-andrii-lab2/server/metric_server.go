package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andrii/apz-pzpi-22-2-ieremenko-andrii/Lab2/pzpi-22-2-ieremenko-andrii-lab2/models"
	"github.com/google/uuid"
)

// MetricServer represents the metric management server
type MetricServer struct {
	metrics  map[uuid.UUID]*models.Metric
	readings map[uuid.UUID][]*models.MetricReading
}

// NewMetricServer creates a new metric server instance
func NewMetricServer() *MetricServer {
	return &MetricServer{
		metrics:  make(map[uuid.UUID]*models.Metric),
		readings: make(map[uuid.UUID][]*models.MetricReading),
	}
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
func (s *MetricServer) CreateMetric(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	metric := &models.Metric{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
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
func (s *MetricServer) ListMetrics(w http.ResponseWriter, r *http.Request) {
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
func (s *MetricServer) GetMetric(w http.ResponseWriter, r *http.Request) {
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
func (s *MetricServer) DeleteMetric(w http.ResponseWriter, r *http.Request) {
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
func (s *MetricServer) AddReading(w http.ResponseWriter, r *http.Request) {
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
func (s *MetricServer) GetReadings(w http.ResponseWriter, r *http.Request) {
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

// Start starts the metric server
func (s *MetricServer) Start(addr string) error {
	// API endpoints
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

	return http.ListenAndServe(addr, nil)
}
