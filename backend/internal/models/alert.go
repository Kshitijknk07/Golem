package models

import "sync"

type Alert struct {
	ID       int    `json:"id"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

var (
	alerts      = []Alert{}
	alertsMutex = &sync.Mutex{}
)

func GetAllAlerts() []Alert {
	alertsMutex.Lock()
	defer alertsMutex.Unlock()
	return alerts
}

func TriggerAlert(alert Alert) {
	alertsMutex.Lock()
	defer alertsMutex.Unlock()
	alert.ID = len(alerts) + 1
	alerts = append(alerts, alert)
}
