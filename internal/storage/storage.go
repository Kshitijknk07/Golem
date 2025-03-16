package storage

import (
	"fmt"
	"sync"
	"time"

	"Golem/internal/metrics"
)

type MetricStorage interface {
	StoreMetrics(metrics metrics.SystemMetrics)
	GetLatestMetrics() metrics.SystemMetrics
	GetMetricsHistory(duration time.Duration) []metrics.SystemMetrics
}

type HealthCheckStorage interface {
	StoreHealthCheckConfig(config metrics.HealthCheckConfig) error
	GetHealthCheckConfig(id string) (metrics.HealthCheckConfig, error)
	GetAllHealthCheckConfigs() []metrics.HealthCheckConfig
	DeleteHealthCheckConfig(id string) error
	StoreHealthCheckResult(result metrics.HealthCheckResult)
	GetHealthCheckResult(id string) (metrics.HealthCheckResult, error)
	GetAllHealthCheckResults() []metrics.HealthCheckResult
	GetHealthCheckHistory(id string, duration time.Duration) []metrics.HealthCheckHistoryEntry
}

type MemoryStorage struct {
	mu             sync.RWMutex
	latestMetrics  metrics.SystemMetrics
	metricsHistory []metrics.SystemMetrics
	maxHistory     int

	healthCheckConfigs map[string]metrics.HealthCheckConfig
	healthCheckResults map[string]metrics.HealthCheckResult
	healthCheckHistory map[string][]metrics.HealthCheckHistoryEntry
	maxCheckHistory    int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		metricsHistory:     make([]metrics.SystemMetrics, 0, 1000),
		maxHistory:         1000,
		healthCheckConfigs: make(map[string]metrics.HealthCheckConfig),
		healthCheckResults: make(map[string]metrics.HealthCheckResult),
		healthCheckHistory: make(map[string][]metrics.HealthCheckHistoryEntry),
		maxCheckHistory:    100,
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

func (s *MemoryStorage) StoreHealthCheckConfig(config metrics.HealthCheckConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.healthCheckConfigs[config.ID] = config
	return nil
}

func (s *MemoryStorage) GetHealthCheckConfig(id string) (metrics.HealthCheckConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	config, exists := s.healthCheckConfigs[id]
	if !exists {
		return metrics.HealthCheckConfig{}, fmt.Errorf("health check config not found: %s", id)
	}

	return config, nil
}

func (s *MemoryStorage) GetAllHealthCheckConfigs() []metrics.HealthCheckConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	configs := make([]metrics.HealthCheckConfig, 0, len(s.healthCheckConfigs))
	for _, config := range s.healthCheckConfigs {
		configs = append(configs, config)
	}

	return configs
}

func (s *MemoryStorage) DeleteHealthCheckConfig(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.healthCheckConfigs[id]; !exists {
		return fmt.Errorf("health check config not found: %s", id)
	}

	delete(s.healthCheckConfigs, id)
	return nil
}

func (s *MemoryStorage) StoreHealthCheckResult(result metrics.HealthCheckResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.healthCheckResults[result.ID] = result

	historyEntry := metrics.HealthCheckHistoryEntry{
		Timestamp:    result.LastChecked,
		Status:       result.Status,
		ResponseTime: result.ResponseTime,
		Message:      result.Message,
	}

	history, exists := s.healthCheckHistory[result.ID]
	if !exists {
		history = make([]metrics.HealthCheckHistoryEntry, 0, s.maxCheckHistory)
	}

	history = append([]metrics.HealthCheckHistoryEntry{historyEntry}, history...)

	if len(history) > s.maxCheckHistory {
		history = history[:s.maxCheckHistory]
	}

	s.healthCheckHistory[result.ID] = history
}

func (s *MemoryStorage) GetHealthCheckResult(id string) (metrics.HealthCheckResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, exists := s.healthCheckResults[id]
	if !exists {
		return metrics.HealthCheckResult{}, fmt.Errorf("health check result not found: %s", id)
	}

	result.History = s.healthCheckHistory[id]

	return result, nil
}

func (s *MemoryStorage) GetAllHealthCheckResults() []metrics.HealthCheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]metrics.HealthCheckResult, 0, len(s.healthCheckResults))
	for id, result := range s.healthCheckResults {
		result.History = s.healthCheckHistory[id]
		results = append(results, result)
	}

	return results
}

func (s *MemoryStorage) GetHealthCheckHistory(id string, duration time.Duration) []metrics.HealthCheckHistoryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history, exists := s.healthCheckHistory[id]
	if !exists {
		return []metrics.HealthCheckHistoryEntry{}
	}

	if duration <= 0 {
		return history
	}

	cutoffTime := time.Now().Add(-duration)
	var result []metrics.HealthCheckHistoryEntry

	for _, entry := range history {
		if entry.Timestamp.After(cutoffTime) {
			result = append(result, entry)
		}
	}

	return result
}
