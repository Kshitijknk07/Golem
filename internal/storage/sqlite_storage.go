package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"Golem/internal/metrics"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	storage := &SQLiteStorage{db: db}
	if err := storage.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %v", err)
	}

	return storage, nil
}

func (s *SQLiteStorage) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			data TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS health_check_configs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			target TEXT NOT NULL,
			interval INTEGER NOT NULL,
			timeout INTEGER NOT NULL,
			enabled BOOLEAN NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS health_check_results (
			id TEXT PRIMARY KEY,
			config_id TEXT NOT NULL,
			status TEXT NOT NULL,
			response_time INTEGER NOT NULL,
			message TEXT,
			last_checked DATETIME NOT NULL,
			FOREIGN KEY (config_id) REFERENCES health_check_configs(id)
		)`,
		`CREATE TABLE IF NOT EXISTS health_check_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			config_id TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			status TEXT NOT NULL,
			response_time INTEGER NOT NULL,
			message TEXT,
			FOREIGN KEY (config_id) REFERENCES health_check_configs(id)
		)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	return nil
}

func (s *SQLiteStorage) StoreMetrics(m metrics.SystemMetrics) error {
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %v", err)
	}

	_, err = s.db.Exec(
		"INSERT INTO metrics (timestamp, data) VALUES (?, ?)",
		m.Timestamp, string(data),
	)
	if err != nil {
		return fmt.Errorf("failed to store metrics: %v", err)
	}

	// Clean up old metrics (keep last 24 hours)
	_, err = s.db.Exec(
		"DELETE FROM metrics WHERE timestamp < datetime('now', '-24 hours')",
	)
	if err != nil {
		return fmt.Errorf("failed to clean up old metrics: %v", err)
	}

	return nil
}

func (s *SQLiteStorage) GetLatestMetrics() (metrics.SystemMetrics, error) {
	var data string
	err := s.db.QueryRow(
		"SELECT data FROM metrics ORDER BY timestamp DESC LIMIT 1",
	).Scan(&data)
	if err == sql.ErrNoRows {
		return metrics.SystemMetrics{}, nil
	}
	if err != nil {
		return metrics.SystemMetrics{}, fmt.Errorf("failed to get latest metrics: %v", err)
	}

	var m metrics.SystemMetrics
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return metrics.SystemMetrics{}, fmt.Errorf("failed to unmarshal metrics: %v", err)
	}

	return m, nil
}

func (s *SQLiteStorage) GetMetricsHistory(duration time.Duration) ([]metrics.SystemMetrics, error) {
	rows, err := s.db.Query(
		"SELECT data FROM metrics WHERE timestamp > datetime('now', ?) ORDER BY timestamp DESC",
		fmt.Sprintf("-%d seconds", int(duration.Seconds())),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics history: %v", err)
	}
	defer rows.Close()

	var result []metrics.SystemMetrics
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan metrics row: %v", err)
		}

		var m metrics.SystemMetrics
		if err := json.Unmarshal([]byte(data), &m); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metrics: %v", err)
		}

		result = append(result, m)
	}

	return result, nil
}

func (s *SQLiteStorage) StoreHealthCheckConfig(config metrics.HealthCheckConfig) error {
	now := time.Now()
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO health_check_configs 
		(id, name, type, target, interval, timeout, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		config.ID, config.Name, config.Type, config.Target,
		config.Interval, config.Timeout, config.Enabled,
		now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to store health check config: %v", err)
	}
	return nil
}

func (s *SQLiteStorage) GetHealthCheckConfig(id string) (metrics.HealthCheckConfig, error) {
	var config metrics.HealthCheckConfig
	err := s.db.QueryRow(
		`SELECT id, name, type, target, interval, timeout, enabled
		FROM health_check_configs WHERE id = ?`,
		id,
	).Scan(
		&config.ID, &config.Name, &config.Type, &config.Target,
		&config.Interval, &config.Timeout, &config.Enabled,
	)
	if err == sql.ErrNoRows {
		return metrics.HealthCheckConfig{}, fmt.Errorf("health check config not found: %s", id)
	}
	if err != nil {
		return metrics.HealthCheckConfig{}, fmt.Errorf("failed to get health check config: %v", err)
	}
	return config, nil
}

func (s *SQLiteStorage) GetAllHealthCheckConfigs() ([]metrics.HealthCheckConfig, error) {
	rows, err := s.db.Query(
		`SELECT id, name, type, target, interval, timeout, enabled
		FROM health_check_configs ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query health check configs: %v", err)
	}
	defer rows.Close()

	var configs []metrics.HealthCheckConfig
	for rows.Next() {
		var config metrics.HealthCheckConfig
		err := rows.Scan(
			&config.ID, &config.Name, &config.Type, &config.Target,
			&config.Interval, &config.Timeout, &config.Enabled,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan health check config: %v", err)
		}
		configs = append(configs, config)
	}
	return configs, nil
}

func (s *SQLiteStorage) DeleteHealthCheckConfig(id string) error {
	_, err := s.db.Exec("DELETE FROM health_check_configs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete health check config: %v", err)
	}
	return nil
}

func (s *SQLiteStorage) StoreHealthCheckResult(result metrics.HealthCheckResult) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Update or insert the latest result
	_, err = tx.Exec(
		`INSERT OR REPLACE INTO health_check_results 
		(id, config_id, status, response_time, message, last_checked)
		VALUES (?, ?, ?, ?, ?, ?)`,
		result.ID, result.Name, result.Status,
		result.ResponseTime, result.Message, result.LastChecked,
	)
	if err != nil {
		return fmt.Errorf("failed to store health check result: %v", err)
	}

	// Add to history
	_, err = tx.Exec(
		`INSERT INTO health_check_history 
		(config_id, timestamp, status, response_time, message)
		VALUES (?, ?, ?, ?, ?)`,
		result.Name, result.LastChecked, result.Status,
		result.ResponseTime, result.Message,
	)
	if err != nil {
		return fmt.Errorf("failed to store health check history: %v", err)
	}

	// Clean up old history entries (keep last 100)
	_, err = tx.Exec(
		`DELETE FROM health_check_history 
		WHERE config_id = ? AND id NOT IN (
			SELECT id FROM health_check_history 
			WHERE config_id = ? 
			ORDER BY timestamp DESC LIMIT 100
		)`,
		result.Name, result.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to clean up old history: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *SQLiteStorage) GetHealthCheckResult(id string) (metrics.HealthCheckResult, error) {
	var result metrics.HealthCheckResult
	err := s.db.QueryRow(
		`SELECT id, config_id, status, response_time, message, last_checked
		FROM health_check_results WHERE id = ?`,
		id,
	).Scan(
		&result.ID, &result.Name, &result.Status,
		&result.ResponseTime, &result.Message, &result.LastChecked,
	)
	if err == sql.ErrNoRows {
		return metrics.HealthCheckResult{}, fmt.Errorf("health check result not found: %s", id)
	}
	if err != nil {
		return metrics.HealthCheckResult{}, fmt.Errorf("failed to get health check result: %v", err)
	}

	// Get history
	rows, err := s.db.Query(
		`SELECT timestamp, status, response_time, message
		FROM health_check_history
		WHERE config_id = ?
		ORDER BY timestamp DESC
		LIMIT 100`,
		result.Name,
	)
	if err != nil {
		return metrics.HealthCheckResult{}, fmt.Errorf("failed to query health check history: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry metrics.HealthCheckHistoryEntry
		err := rows.Scan(
			&entry.Timestamp, &entry.Status,
			&entry.ResponseTime, &entry.Message,
		)
		if err != nil {
			return metrics.HealthCheckResult{}, fmt.Errorf("failed to scan history entry: %v", err)
		}
		result.History = append(result.History, entry)
	}

	return result, nil
}

func (s *SQLiteStorage) GetAllHealthCheckResults() ([]metrics.HealthCheckResult, error) {
	rows, err := s.db.Query(
		`SELECT id, config_id, status, response_time, message, last_checked
		FROM health_check_results
		ORDER BY last_checked DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query health check results: %v", err)
	}
	defer rows.Close()

	var results []metrics.HealthCheckResult
	for rows.Next() {
		var result metrics.HealthCheckResult
		err := rows.Scan(
			&result.ID, &result.Name, &result.Status,
			&result.ResponseTime, &result.Message, &result.LastChecked,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan health check result: %v", err)
		}

		// Get history for this result
		historyRows, err := s.db.Query(
			`SELECT timestamp, status, response_time, message
			FROM health_check_history
			WHERE config_id = ?
			ORDER BY timestamp DESC
			LIMIT 100`,
			result.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query health check history: %v", err)
		}

		for historyRows.Next() {
			var entry metrics.HealthCheckHistoryEntry
			err := historyRows.Scan(
				&entry.Timestamp, &entry.Status,
				&entry.ResponseTime, &entry.Message,
			)
			if err != nil {
				historyRows.Close()
				return nil, fmt.Errorf("failed to scan history entry: %v", err)
			}
			result.History = append(result.History, entry)
		}
		historyRows.Close()

		results = append(results, result)
	}

	return results, nil
}

func (s *SQLiteStorage) GetHealthCheckHistory(id string, duration time.Duration) ([]metrics.HealthCheckHistoryEntry, error) {
	query := `SELECT timestamp, status, response_time, message
		FROM health_check_history
		WHERE config_id = ?`
	args := []interface{}{id}

	if duration > 0 {
		query += ` AND timestamp > datetime('now', ?)`
		args = append(args, fmt.Sprintf("-%d seconds", int(duration.Seconds())))
	}

	query += ` ORDER BY timestamp DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query health check history: %v", err)
	}
	defer rows.Close()

	var history []metrics.HealthCheckHistoryEntry
	for rows.Next() {
		var entry metrics.HealthCheckHistoryEntry
		err := rows.Scan(
			&entry.Timestamp, &entry.Status,
			&entry.ResponseTime, &entry.Message,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history entry: %v", err)
		}
		history = append(history, entry)
	}

	return history, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
