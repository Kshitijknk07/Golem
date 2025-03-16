package metrics

import (
	"time"
)

type HealthCheckType string

const (
	HTTPCheck     HealthCheckType = "http"
	TCPCheck      HealthCheckType = "tcp"
	DatabaseCheck HealthCheckType = "database"
	APICheck      HealthCheckType = "api"
)

type HealthCheckStatus string

const (
	StatusUp      HealthCheckStatus = "up"
	StatusDown    HealthCheckStatus = "down"
	StatusWarning HealthCheckStatus = "warning"
	StatusUnknown HealthCheckStatus = "unknown"
)

type HealthCheckConfig struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Type       HealthCheckType   `json:"type"`
	Target     string            `json:"target"`
	Interval   time.Duration     `json:"interval"`
	Timeout    time.Duration     `json:"timeout"`
	Method     string            `json:"method,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
	ExpectCode int               `json:"expect_code,omitempty"`
	ExpectBody string            `json:"expect_body,omitempty"`
	Enabled    bool              `json:"enabled"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type HealthCheckResult struct {
	ID           string                    `json:"id"`
	Name         string                    `json:"name"`
	Type         HealthCheckType           `json:"type"`
	Target       string                    `json:"target"`
	Status       HealthCheckStatus         `json:"status"`
	ResponseTime time.Duration             `json:"response_time"`
	Message      string                    `json:"message,omitempty"`
	LastChecked  time.Time                 `json:"last_checked"`
	History      []HealthCheckHistoryEntry `json:"history,omitempty"`
}

type HealthCheckHistoryEntry struct {
	Timestamp    time.Time         `json:"timestamp"`
	Status       HealthCheckStatus `json:"status"`
	ResponseTime time.Duration     `json:"response_time"`
	Message      string            `json:"message,omitempty"`
}

type HealthCheckMetrics struct {
	Checks []HealthCheckResult `json:"checks"`
}
