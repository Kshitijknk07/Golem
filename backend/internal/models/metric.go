package models

import (
	"golem/internal/db"
	"time"
)

type Metric struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Value     float64     `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func GetAllMetrics() ([]Metric, error) {
	rows, err := db.DB.Query("SELECT id, name, value, timestamp FROM metrics ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.ID, &m.Name, &m.Value, &m.Timestamp); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, rows.Err()
}

func SaveMetric(metric Metric) error {
	query := "INSERT INTO metrics (name, value) VALUES (?, ?)"
	result, err := db.DB.Exec(query, metric.Name, metric.Value)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	metric.ID = int(id)
	return nil
}

func UpdateMetric(metric Metric) error {
	query := "UPDATE metrics SET name = ?, value = ? WHERE id = ?"
	_, err := db.DB.Exec(query, metric.Name, metric.Value, metric.ID)
	return err
}
