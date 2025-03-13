package models

import (
	"database/sql"
	"golem/internal/db"
	"time"
)

type Alert struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
}

func GetAllAlerts() ([]Alert, error) {
	rows, err := db.DB.Query("SELECT id, message, severity, timestamp FROM alerts ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var a Alert
		if err := rows.Scan(&a.ID, &a.Message, &a.Severity, &a.Timestamp); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func TriggerAlert(alert Alert) error {
	query := "INSERT INTO alerts (message, severity) VALUES (?, ?)"
	result, err := db.DB.Exec(query, alert.Message, alert.Severity)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	alert.ID = int(id)
	return nil
}
