package models

import (
	"time"

	"github.com/google/uuid"
)

// Metric represents a household metric
type Metric struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"` // единица измерения (кВт, м³, и т.д.)
	RoomID      uuid.UUID `json:"room_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MetricReading represents a single reading of a metric
type MetricReading struct {
	ID        uuid.UUID `json:"id"`
	MetricID  uuid.UUID `json:"metric_id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateMetricRequest represents the request to create a new metric
type CreateMetricRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Unit        string    `json:"unit" binding:"required"`
	RoomID      uuid.UUID `json:"room_id"`
}

// AddReadingRequest represents the request to add a new reading
type AddReadingRequest struct {
	Value     float64   `json:"value" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
}

// MetricWithReadings represents a metric with its readings
type MetricWithReadings struct {
	Metric   Metric          `json:"metric"`
	Readings []MetricReading `json:"readings"`
}

// MetricListResponse represents the response for listing metrics
type MetricListResponse struct {
	Metrics []Metric `json:"metrics"`
	Total   int      `json:"total"`
}

// ReadingListResponse represents the response for listing readings
type ReadingListResponse struct {
	Readings []MetricReading `json:"readings"`
	Total    int             `json:"total"`
}

// CorrelationRequest represents a request to calculate correlation between metrics
type CorrelationRequest struct {
	Metric1ID uuid.UUID `json:"metric1Id"`
	Metric2ID uuid.UUID `json:"metric2Id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// CorrelationResponse represents the correlation calculation result
type CorrelationResponse struct {
	Metric1Name string    `json:"metric1Name"`
	Metric2Name string    `json:"metric2Name"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Correlation float64   `json:"correlation"`
	Message     string    `json:"message"`
}
