package services

import (
	"golem/internal/models"
	"time"
)

var alertsStore = make(map[string]models.Alert)

func TriggerAlert(metric models.Metric) (*models.Alert, error) {
	var alert models.Alert

	if metric.Type == "CPU" && metric.Value > 90 {
		alert = models.Alert{
			ID:        generateID(),
			MetricID:  metric.ID,
			Threshold: 90,
			Severity:  "Critical",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		alertsStore[alert.ID] = alert
	}

	return &alert, nil
}

func generateID() string {
	return time.Now().Format("20060102150405")
}

func GetAlerts() []models.Alert {
	var result []models.Alert
	for _, alert := range alertsStore {
		result = append(result, alert)
	}
	return result
}
