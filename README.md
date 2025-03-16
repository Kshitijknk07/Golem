# Golem - Distributed System Monitoring Platform

Golem is a monitoring solution for distributed systems, providing real-time metrics collection, health checking, and visualization capabilities. Built with Go, it offers a lightweight yet powerful approach to system monitoring with minimal resource footprint.

## 🚀 Features

### 📊 System Metrics Monitoring
- **CPU Metrics**: Total usage, per-core usage, load averages
- **Memory Metrics**: RAM usage, swap usage, usage percentages
- **Disk Metrics**: Partition usage, I/O statistics
- **Network Metrics**: Interface statistics, bandwidth usage
- **Process Monitoring**: Top processes by CPU and memory usage
- **System Uptime**: Boot time and system uptime tracking
- **Historical Data**: Time-series storage of all metrics

### 🩺 Health Check System
- **HTTP/HTTPS Endpoint Checks**: Status code validation, response time monitoring
- **TCP Port Availability**: Connection testing to network services
- **Database Connectivity**: MySQL, PostgreSQL, SQLite connection validation
- **API Endpoint Monitoring**: JSON response validation
- **Custom Checks**: Plugin system for specialized checks
- **Configurable Parameters**: Intervals, timeouts, expected responses

### 🌐 Web Dashboard
- **Real-time Metrics Visualization**: Live updates without page refresh
- **Health Check Status Overview**: At-a-glance service health status
- **Responsive Design**: Mobile and desktop-friendly interface
- **Interactive Controls**: Add, edit, and delete health checks from UI

### 🔌 Extensible Architecture
- **Plugin System**: Extend functionality with custom health checks
- **Storage Backends**: Memory, MySQL
- **REST API**: Full programmatic access to all features

---

## 🏗️ Architecture

### System Architecture Diagram
![dg1](https://github.com/user-attachments/assets/39289110-2712-40df-ae59-28f42abda4e3)

### Component Flow Diagram
![dg2](https://github.com/user-attachments/assets/78cd7c53-d1e9-42ad-be59-6f67b9b10aa3)

## 🏁 Getting Started

### Prerequisites

- **Go 1.24+**
- **MySQL (optional for persistent storage)**
- **Modern web browser** for dashboard access

---

### Installation

#### From Source
```bash
# Clone the repository
git clone https://github.com/yourusername/Golem.git
cd Golem

# Build the executable
go build -o golem cmd/golem/main.go

# Run Golem
./golem

# Access the web dashboard at: http://localhost:8080
```

#### Using Pre-built Binaries
```bash
# Download the latest release for your platform

# Windows
curl -LO https://github.com/yourusername/Golem/releases/latest/download/golem-windows-amd64.zip
unzip golem-windows-amd64.zip

# Run Golem
golem.exe

# Access the web dashboard at: http://localhost:8080
```

---

### Configuration

Create a `config.yaml` file in the same directory as the Golem executable:

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  static_dir: "web/static"
  log_level: "info"  # debug, info, warn, error

storage:
  type: "memory"  # Options: memory | mysql
  retention_period: "720h"  # 30 days

  # MySQL configuration (if type is "mysql")
  mysql:
    host: "localhost"
    port: 3306
    database: "golem"
    username: "golem_user"
    password: "password"

health_checks:
  default_interval: "60s"
  default_timeout: "10s"
  history_retention: "168h"  # 7 days

metrics:
  collection_interval: "5s"
  process_limit: 20  # Number of top processes to track
```

---

## 💻 Usage

### Web Dashboard

The Golem web dashboard provides a comprehensive view of your system's health and performance metrics. Access it by navigating to [http://localhost:8080](http://localhost:8080) in your web browser after starting the Golem service.

#### Dashboard Sections

- **System Overview:** At-a-glance view of key system metrics:
  - CPU Usage
  - Memory Usage
  - Disk Usage
  - System Uptime
- **CPU Details:** Per-core usage and load averages
- **Memory Details:** RAM and swap usage statistics
- **Disk Details:** Usage by partition and I/O statistics
- **Network Details:** Interface statistics and bandwidth usage
- **Top Processes:** Resource usage by process
- **Health Checks & Service Monitoring:** Status of monitored services:
  - Add, edit, and delete health checks
  - View health check history and response times

---

### API Reference

Golem provides a comprehensive REST API for programmatic access to all features.

#### Metrics Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| GET | `/api/metrics` | Fetch the latest system metrics |
| GET | `/api/metrics/history?duration=1h` | Fetch historical system metrics for a specified duration |
| GET | `/api/metrics/cpu` | Fetch only CPU metrics |
| GET | `/api/metrics/memory` | Fetch only memory metrics |
| GET | `/api/metrics/disk` | Fetch only disk metrics |
| GET | `/api/metrics/network` | Fetch only network metrics |
| GET | `/api/metrics/processes` | Fetch only process metrics |

#### Health Check Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| GET | `/api/health-checks` | List all health checks |
| POST | `/api/health-checks` | Create a new health check |
| GET | `/api/health-checks/{id}` | Get details of a specific health check |
| PUT | `/api/health-checks/{id}` | Update an existing health check |
| DELETE | `/api/health-checks/{id}` | Delete a health check |
| GET | `/api/health-checks/{id}/history?duration=24h` | Fetch the history of a health check |
| POST | `/api/health-checks/{id}/run` | Manually trigger a health check |
| PUT | `/api/health-checks/{id}/enable` | Enable a health check |
| PUT | `/api/health-checks/{id}/disable` | Disable a health check |

## System Metrics

### CPU Metrics
```json
{
  "total_usage": 25.5,
  "per_core_usage": {
    "0": 30.2,
    "1": 22.8,
    "2": 18.5,
    "3": 30.5
  },
  "load_average": [1.25, 1.15, 0.95]
}
```

### Memory Metrics
```json
{
  "total": 16777216000,
  "used": 8388608000,
  "free": 8388608000,
  "used_percent": 50.0,
  "swap_total": 4294967296,
  "swap_used": 1073741824,
  "swap_free": 3221225472
}
```

### Disk Metrics
```json
{
  "partitions": [
    {
      "device": "C:",
      "mountpoint": "C:",
      "total": 256060514304,
      "used": 102424205696,
      "free": 153636308608,
      "used_percent": 40.0
    }
  ],
  "io_counters": {
    "C:": {
      "read_count": 123456,
      "write_count": 78901,
      "read_bytes": 1073741824,
      "write_bytes": 536870912,
      "read_time": 1500,
      "write_time": 1200
    }
  }
}
```

### Network Metrics
```json
{
  "interfaces": {
    "eth0": {
      "bytes_sent": 1073741824,
      "bytes_recv": 2147483648,
      "packets_sent": 8192,
      "packets_recv": 16384,
      "errin": 0,
      "errout": 0,
      "dropin": 0,
      "dropout": 0
    }
  }
}
```

### Process Metrics
```json
[
  {
    "pid": 1234,
    "name": "chrome.exe",
    "username": "user",
    "cpu_percent": 5.2,
    "memory_used": 1073741824,
    "status": "running",
    "create_time": 1620000000,
    "num_threads": 32,
    "io_counters": {
      "read_count": 1000,
      "write_count": 500,
      "read_bytes": 10485760,
      "write_bytes": 5242880
    }
  }
]
```

## Health Checks

### Health Check Types

#### HTTP/HTTPS Check
```json
{
  "name": "Company Website",
  "type": "http",
  "target": "https://example.com",
  "interval": "60s",
  "timeout": "10s",
  "method": "GET",
  "headers": {
    "User-Agent": "Golem-Monitor/1.0"
  },
  "expectCode": 200,
  "expectBody": "Welcome"
}
```

#### TCP Check
```json
{
  "name": "Database Server",
  "type": "tcp",
  "target": "db.example.com:5432",
  "interval": "30s",
  "timeout": "5s"
}
```

#### Database Check
```json
{
  "name": "MySQL Database",
  "type": "database",
  "target": "mysql://user:password@localhost:3306/dbname",
  "interval": "60s",
  "timeout": "10s",
  "query": "SELECT 1"
}
```

#### API Check
```json
{
  "name": "User API",
  "type": "api",
  "target": "https://api.example.com/health",
  "interval": "60s",
  "timeout": "10s",
  "method": "GET",
  "headers": {
    "Authorization": "Bearer token123"
  },
  "expectCode": 200,
  "expectJson": {
    "status": "healthy"
  }
}
```

### Health Check Results
```json
{
  "id": "chk_123456",
  "name": "Example Website",
  "type": "http",
  "target": "https://example.com",
  "status": "up",
  "response_time": 235000000,
  "message": "HTTP 200 OK",
  "last_checked": "2023-05-01T12:34:56Z"
}
```

## Storage Backends

### Memory Storage
Default in-memory storage suitable for testing and short-term monitoring. Data is lost when the service restarts.

### MySQL Storage
Persistent storage using MySQL database:

```yaml
storage:
  type: "mysql"
  mysql:
    host: "localhost"
    port: 3306
    database: "golem"
    username: "golem_user"
    password: "password"
```
## System Metrics

### CPU Metrics
```json
{
  "total_usage": 25.5,
  "per_core_usage": {
    "0": 30.2,
    "1": 22.8,
    "2": 18.5,
    "3": 30.5
  },
  "load_average": [1.25, 1.15, 0.95]
}
```

### Memory Metrics
```json
{
  "total": "16GB",
  "used": "8GB",
  "free": "8GB",
  "used_percent": 50.0,
  "swap_total": "4GB",
  "swap_used": "1GB",
  "swap_free": "3GB"
}
```

### Disk Metrics
```json
{
  "partitions": [
    {
      "device": "C:",
      "mountpoint": "C:",
      "total": "256GB",
      "used": "102GB",
      "free": "154GB",
      "used_percent": 40.0
    }
  ],
  "io_counters": {
    "C:": {
      "read_count": 123456,
      "write_count": 78901,
      "read_bytes": "1GB",
      "write_bytes": "512MB",
      "read_time": 1500,
      "write_time": 1200
    }
  }
}
```

### Network Metrics
```json
{
  "interfaces": {
    "eth0": {
      "bytes_sent": "1GB",
      "bytes_recv": "2GB",
      "packets_sent": 8192,
      "packets_recv": 16384,
      "errin": 0,
      "errout": 0,
      "dropin": 0,
      "dropout": 0
    }
  }
}
```

### Process Metrics
```json
[
  {
    "pid": 1234,
    "name": "chrome.exe",
    "username": "user",
    "cpu_percent": 5.2,
    "memory_used": "1GB",
    "status": "running",
    "create_time": 1620000000,
    "num_threads": 32,
    "io_counters": {
      "read_count": 1000,
      "write_count": 500,
      "read_bytes": "10MB",
      "write_bytes": "5MB"
    }
  }
]
```

## Health Checks

### Health Check Types

#### HTTP/HTTPS Check
```json
{
  "name": "Company Website",
  "type": "http",
  "target": "https://example.com",
  "interval": "60s",
  "timeout": "10s",
  "method": "GET",
  "headers": {
    "User-Agent": "Golem-Monitor/1.0"
  },
  "expectCode": 200,
  "expectBody": "Welcome"
}
```

#### TCP Check
```json
{
  "name": "Database Server",
  "type": "tcp",
  "target": "db.example.com:5432",
  "interval": "30s",
  "timeout": "5s"
}
```

#### Database Check
```json
{
  "name": "MySQL Database",
  "type": "database",
  "target": "mysql://user:password@localhost:3306/dbname",
  "interval": "60s",
  "timeout": "10s",
  "query": "SELECT 1"
}
```

#### API Check
```json
{
  "name": "User API",
  "type": "api",
  "target": "https://api.example.com/health",
  "interval": "60s",
  "timeout": "10s",
  "method": "GET",
  "headers": {
    "Authorization": "Bearer token123"
  },
  "expectCode": 200,
  "expectJson": {
    "status": "healthy"
  }
}
```

### Health Check Results
```json
{
  "id": "chk_123456",
  "name": "Example Website",
  "type": "http",
  "target": "https://example.com",
  "status": "up",
  "response_time": 235000000,
  "message": "HTTP 200 OK",
  "last_checked": "2023-05-01T12:34:56Z"
}
```

## Storage Backends

### Memory Storage
Default in-memory storage suitable for testing and short-term monitoring. Data is lost when the service restarts.

### MySQL Storage
Persistent storage using MySQL database:

```yaml
storage:
  type: "mysql"
  mysql:
    host: "localhost"
    port: 3306
    database: "golem"
    username: "golem_user"
    password: "password"
```

## Plugin System

Golem features an extensible plugin system for custom health checks.

### Plugin Interface
```go
type CheckPlugin interface {
    Name() string
    Type() metrics.HealthCheckType
    Description() string
    Execute(ctx context.Context, target string, timeout time.Duration) (metrics.HealthCheckStatus, string, time.Duration)
    ValidateConfig(config map[string]interface{}) error
}
```

### Creating Custom Plugins
To create a custom health check plugin:
1. Implement the `CheckPlugin` interface.
2. Register your plugin with the plugin registry.
3. Configure health checks to use your plugin.

### Example Custom Plugin
```go
package myplugin

import (
    "Golem/internal/metrics"
    "Golem/internal/plugin"
    "context"
    "time"
)

type MyCustomCheck struct{}

func (c *MyCustomCheck) Name() string {
    return "my-custom-check"
}

func (c *MyCustomCheck) Type() metrics.HealthCheckType {
    return "custom"
}

func (c *MyCustomCheck) Description() string {
    return "Custom health check for specialized service monitoring"
}

func (c *MyCustomCheck) Execute(ctx context.Context, target string, timeout time.Duration) (metrics.HealthCheckStatus, string, time.Duration) {
    return metrics.StatusUp, "Service is healthy", 150 * time.Millisecond
}

func (c *MyCustomCheck) ValidateConfig(config map[string]interface{}) error {
    return nil
}

// Register the plugin
func init() {
    registry := plugin.NewRegistry()
    registry.Register(&MyCustomCheck{})
}
```

## Data Models

### System Metrics
```go
 type SystemMetrics struct {
     Timestamp   time.Time              `json:"timestamp"`
     CPU         CPUMetrics             `json:"cpu"`
     Memory      MemoryMetrics          `json:"memory"`
     Disk        DiskMetrics            `json:"disk"`
     Network     NetworkMetrics         `json:"network"`
     Process     []ProcessMetrics       `json:"processes"`
     Uptime      UptimeMetrics          `json:"uptime"`
     HealthCheck HealthCheckMetrics     `json:"health_checks"`
 }
```

### Health Check Configuration
```go
 type HealthCheckConfig struct {
     ID         string                `json:"id"`
     Name       string                `json:"name"`
     Type       HealthCheckType       `json:"type"`
     Target     string                `json:"target"`
     Interval   time.Duration         `json:"interval"`
     Timeout    time.Duration         `json:"timeout"`
     Method     string                `json:"method,omitempty"`
     Headers    map[string]string     `json:"headers,omitempty"`
     Body       string                `json:"body,omitempty"`
     ExpectCode int                   `json:"expect_code,omitempty"`
     ExpectBody string                `json:"expect_body,omitempty"`
     PluginName string                `json:"plugin_name,omitempty"`
     Enabled    bool                  `json:"enabled"`
     CreatedAt  time.Time             `json:"created_at"`
     UpdatedAt  time.Time             `json:"updated_at"`
 }
```

## 🔄 Internal Architecture

### Core Components

#### Collector
- Gathers system metrics at regular intervals
- Collects:
  - CPU, memory, disk, network, and process metrics
- Uses the `gopsutil` library for cross-platform compatibility

#### Health Check Collector
- Manages and executes health checks
- Supports:
  - HTTP/HTTPS, TCP, database, and API checks
- Plugin support for custom checks

#### Storage
- Manages data persistence
- Supports:
  - In-memory storage for development/testing
  - Database backends for production use

#### API Server
- Provides REST API endpoints for:
  - Metrics retrieval
  - Health check management
  - Historical data access

#### Web Dashboard
- User interface for monitoring
- Features:
  - Real-time metrics visualization
  - Health check management
  - Responsive design

## 📊 Sequence Diagrams

### Metrics Collection Flow
![dg3](https://github.com/user-attachments/assets/ae8b37a7-bfce-4236-b2c8-bb4238721a76)

### Health Check Execution Flow
![dg4](https://github.com/user-attachments/assets/235cbf04-9a43-468a-a46b-80b838a3e10f)


## 🙌 Acknowledgements

- `gopsutil` - Cross-platform system information library
- `gorilla/mux` - HTTP router and URL matcher
- `go-sql-driver/mysql` - MySQL driver for Go
- `lib/pq` - PostgreSQL driver for Go
- `mattn/go-sqlite3` - SQLite driver for Go



