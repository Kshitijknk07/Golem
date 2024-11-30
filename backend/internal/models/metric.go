package models

import "sync"

type Metric struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

var (
	metrics      = []Metric{}
	metricsMutex = &sync.Mutex{}
)

func GetAllMetrics() []Metric {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	return metrics
}

func SaveMetric(metric Metric) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	metric.ID = len(metrics) + 1
	metrics = append(metrics, metric)
}

func UpdateMetric(updatedMetric Metric) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	for i, metric := range metrics {
		if metric.ID == updatedMetric.ID {
			metrics[i] = updatedMetric
			return
		}
	}
}
