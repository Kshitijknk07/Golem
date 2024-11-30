package service

import "golem/internal/models"

func CreateAlert(message string, severity string) models.Alert {
	alert := models.Alert{
		Message:  message,
		Severity: severity,
	}
	models.TriggerAlert(alert)
	return alert
}
