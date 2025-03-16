package collector

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"Golem/internal/metrics"
	"Golem/internal/plugin"
	"Golem/internal/storage"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type HealthCheckCollector struct {
	storage        storage.HealthCheckStorage
	client         *http.Client
	checkInterval  time.Duration
	mu             sync.RWMutex
	checks         map[string]metrics.HealthCheckConfig
	results        map[string]metrics.HealthCheckResult
	pluginRegistry *plugin.Registry
}

func NewHealthCheckCollector(storage storage.HealthCheckStorage) *HealthCheckCollector {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	registry := plugin.NewRegistry()

	return &HealthCheckCollector{
		storage:        storage,
		client:         client,
		checkInterval:  1 * time.Minute,
		checks:         make(map[string]metrics.HealthCheckConfig),
		results:        make(map[string]metrics.HealthCheckResult),
		pluginRegistry: registry,
	}
}

func (c *HealthCheckCollector) Start(ctx context.Context) {
	configs := c.storage.GetAllHealthCheckConfigs()
	for _, config := range configs {
		c.checks[config.ID] = config
	}

	for _, check := range c.checks {
		if check.Enabled {
			go c.runCheck(ctx, check)
		}
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			configs := c.storage.GetAllHealthCheckConfigs()

			c.mu.Lock()
			for _, config := range configs {
				existingCheck, exists := c.checks[config.ID]
				if !exists || existingCheck.UpdatedAt.Before(config.UpdatedAt) {
					c.checks[config.ID] = config
					if config.Enabled {
						go c.runCheck(ctx, config)
					}
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *HealthCheckCollector) runCheck(ctx context.Context, check metrics.HealthCheckConfig) {
	ticker := time.NewTicker(check.Interval)
	defer ticker.Stop()

	c.performCheck(check)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.performCheck(check)
		}
	}
}

func (c *HealthCheckCollector) performCheck(check metrics.HealthCheckConfig) {
	var result metrics.HealthCheckResult
	result.ID = check.ID
	result.Name = check.Name
	result.Type = check.Type
	result.Target = check.Target
	result.LastChecked = time.Now()

	startTime := time.Now()

	switch check.Type {
	case metrics.HTTPCheck:
		c.performHTTPCheck(&result, check)
	case metrics.TCPCheck:
		c.performTCPCheck(&result, check)
	case metrics.DatabaseCheck:
		c.performDatabaseCheck(&result, check)
	case metrics.APICheck:
		c.performAPICheck(&result, check)
	case "plugin":
		c.performPluginCheck(&result, check)
	default:
		result.Status = metrics.StatusUnknown
		result.Message = "Unknown check type"
	}

	result.ResponseTime = time.Since(startTime)

	c.mu.Lock()
	c.results[check.ID] = result
	c.mu.Unlock()

	c.storage.StoreHealthCheckResult(result)
}

func (c *HealthCheckCollector) performHTTPCheck(result *metrics.HealthCheckResult, check metrics.HealthCheckConfig) {
	method := "GET"
	if check.Method != "" {
		method = check.Method
	}

	req, err := http.NewRequest(method, check.Target, strings.NewReader(check.Body))
	if err != nil {
		result.Status = metrics.StatusDown
		result.Message = fmt.Sprintf("Error creating request: %v", err)
		return
	}

	for key, value := range check.Headers {
		req.Header.Add(key, value)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Golem-Monitoring/1.0")
	}

	client := *c.client
	if check.Timeout > 0 {
		client.Timeout = check.Timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		result.Status = metrics.StatusDown
		result.Message = fmt.Sprintf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	if check.ExpectCode > 0 && resp.StatusCode != check.ExpectCode {
		result.Status = metrics.StatusWarning
		result.Message = fmt.Sprintf("Expected status code %d, got %d", check.ExpectCode, resp.StatusCode)
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = metrics.StatusUp
		result.Message = fmt.Sprintf("HTTP %d %s", resp.StatusCode, resp.Status)
	} else {
		result.Status = metrics.StatusDown
		result.Message = fmt.Sprintf("HTTP %d %s", resp.StatusCode, resp.Status)
	}
}

func (c *HealthCheckCollector) performTCPCheck(result *metrics.HealthCheckResult, check metrics.HealthCheckConfig) {
	timeout := 5 * time.Second
	if check.Timeout > 0 {
		timeout = check.Timeout
	}

	conn, err := net.DialTimeout("tcp", check.Target, timeout)
	if err != nil {
		result.Status = metrics.StatusDown
		result.Message = fmt.Sprintf("Connection failed: %v", err)
		return
	}
	defer conn.Close()

	result.Status = metrics.StatusUp
	result.Message = "TCP connection successful"
}

func (c *HealthCheckCollector) performDatabaseCheck(result *metrics.HealthCheckResult, check metrics.HealthCheckConfig) {
	var driverName string
	if strings.HasPrefix(check.Target, "postgres://") {
		driverName = "postgres"
	} else if strings.HasPrefix(check.Target, "mysql://") {
		driverName = "mysql"
	} else if strings.HasPrefix(check.Target, "sqlite://") {
		driverName = "sqlite3"
		check.Target = strings.TrimPrefix(check.Target, "sqlite://")
	} else {
		result.Status = metrics.StatusDown
		result.Message = "Unsupported database type"
		return
	}

	db, err := sql.Open(driverName, check.Target)
	if err != nil {
		result.Status = metrics.StatusDown
		result.Message = fmt.Sprintf("Database connection error: %v", err)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), check.Timeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		result.Status = metrics.StatusDown
		result.Message = fmt.Sprintf("Database ping failed: %v", err)
		return
	}

	result.Status = metrics.StatusUp
	result.Message = "Database connection successful"
}

func (c *HealthCheckCollector) performAPICheck(result *metrics.HealthCheckResult, check metrics.HealthCheckConfig) {
	c.performHTTPCheck(result, check)
}

func (c *HealthCheckCollector) performPluginCheck(result *metrics.HealthCheckResult, check metrics.HealthCheckConfig) {
	pluginName := check.PluginName
	if pluginName == "" {
		result.Status = metrics.StatusUnknown
		result.Message = "No plugin specified"
		return
	}

	plugin, exists := c.pluginRegistry.Get(pluginName)
	if !exists {
		result.Status = metrics.StatusUnknown
		result.Message = fmt.Sprintf("Plugin '%s' not found", pluginName)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), check.Timeout)
	defer cancel()

	status, message, responseTime := plugin.Execute(ctx, check.Target, check.Timeout)

	result.Status = status
	result.Message = message
	result.ResponseTime = responseTime
}

func (c *HealthCheckCollector) GetHealthCheckResults() metrics.HealthCheckMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	results := metrics.HealthCheckMetrics{
		Checks: make([]metrics.HealthCheckResult, 0, len(c.results)),
	}

	for _, result := range c.results {
		results.Checks = append(results.Checks, result)
	}

	return results
}

func (c *HealthCheckCollector) AddHealthCheck(config metrics.HealthCheckConfig) error {
	if config.Name == "" {
		return fmt.Errorf("health check name cannot be empty")
	}
	if config.Target == "" {
		return fmt.Errorf("health check target cannot be empty")
	}
	if config.Interval == 0 {
		config.Interval = 60 * time.Second
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	if config.ID == "" {
		config.ID = fmt.Sprintf("check_%d", time.Now().UnixNano())
	}

	config.Enabled = true
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	err := c.storage.StoreHealthCheckConfig(config)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.checks[config.ID] = config
	c.mu.Unlock()

	go c.runCheck(context.Background(), config)

	return nil
}

func (c *HealthCheckCollector) UpdateHealthCheck(config metrics.HealthCheckConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.checks[config.ID]
	if !exists {
		return fmt.Errorf("health check with ID %s not found", config.ID)
	}

	config.UpdatedAt = time.Now()

	err := c.storage.StoreHealthCheckConfig(config)
	if err != nil {
		return err
	}

	c.checks[config.ID] = config
	return nil
}

func (c *HealthCheckCollector) DeleteHealthCheck(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.checks[id]; !exists {
		return fmt.Errorf("health check with ID %s not found", id)
	}

	err := c.storage.DeleteHealthCheckConfig(id)
	if err != nil {
		return err
	}

	delete(c.checks, id)
	delete(c.results, id)

	return nil
}

func (c *HealthCheckCollector) GetHealthCheckByID(id string) (metrics.HealthCheckResult, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result, exists := c.results[id]
	if !exists {
		return metrics.HealthCheckResult{}, fmt.Errorf("health check result with ID %s not found", id)
	}

	return result, nil
}
