package storage

import (
	"sync"
	"time"

	"Golem/internal/metrics"
)

type MetricStorage interface {
	StoreMetrics(metrics metrics.SystemMetrics)
	GetLatestMetrics() metrics.SystemMetrics
	GetMetricsHistory(duration time.Duration) []metrics.SystemMetrics
}

type MemoryStorage struct {
	mu             sync.RWMutex
	latestMetrics  metrics.SystemMetrics
	metricsHistory []metrics.SystemMetrics
	maxHistory     int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		metricsHistory: make([]metrics.SystemMetrics, 0, 1000),
		maxHistory:     1000,
	}
}

func (s *MemoryStorage) StoreMetrics(m metrics.SystemMetrics) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.latestMetrics = m

	s.metricsHistory = append(s.metricsHistory, m)
	if len(s.metricsHistory) > s.maxHistory {
		s.metricsHistory = s.metricsHistory[1:]
	}
}

func (s *MemoryStorage) GetLatestMetrics() metrics.SystemMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.latestMetrics
}

func (s *MemoryStorage) GetMetricsHistory(duration time.Duration) []metrics.SystemMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.metricsHistory) == 0 {
		return []metrics.SystemMetrics{}
	}

	cutoffTime := time.Now().Add(-duration)
	var result []metrics.SystemMetrics

	for i := len(s.metricsHistory) - 1; i >= 0; i-- {
		if s.metricsHistory[i].Timestamp.Before(cutoffTime) {
			break
		}
		result = append([]metrics.SystemMetrics{s.metricsHistory[i]}, result...)
	}

	return result
}
