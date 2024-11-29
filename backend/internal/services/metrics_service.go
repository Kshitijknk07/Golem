package services

import (
	"errors"
	"golem/internal/models"
	"time"
)

var metricsStore = make(map[string]models.Metric)

func GetMetrics() ([]models.Metric, error) {
	var result []models.Metric
	for _, metric := range metricsStore {
		result = append(result, metric)
	}
	return result, nil
}

func SaveMetric(newMetric models.Metric) error {
	newMetric.Time = time.Now().Format(time.RFC3339)
	metricsStore[newMetric.ID] = newMetric
	return nil
}

func UpdateMetric(id string, updatedMetric models.Metric) error {
	if _, exists := metricsStore[id]; !exists {
		return errors.New("metric not found")
	}
	metricsStore[id] = updatedMetric
	return nil
}
