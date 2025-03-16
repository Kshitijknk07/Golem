# Golem - Distributed System Monitoring Platform

![Go Version](https://img.shields.io/badge/go-1.24%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)

Golem is an monitoring solution for distributed systems, providing real-time metrics collection, health checking, and visualization capabilities.

## Features

- 📊 **System Metrics Monitoring**
  - CPU/Memory/Disk Usage
  - Network I/O Statistics
  - Process Monitoring
  - Historical Data Storage
- 🩺 **Health Check System**
  - HTTP/HTTPS Endpoint Checks
  - TCP Port Availability
  - Database Connectivity (MySQL, PostgreSQL, SQLite)
  - Customizable Check Intervals
- 🌐 **Web Dashboard**
  - Real-time Metrics Visualization
  - Health Check Status Overview
  - Historical Trend Analysis
- 🔌 **Extensible Architecture**
  - Plugin System for Custom Checks
  - Multiple Storage Backends
  - REST API for Integration

## Getting Started

### Prerequisites

- Go 1.24+
- MySQL/PostgreSQL (optional for persistent storage)

### Installation

```bash
# Clone repository
git clone https://github.com/yourusername/Golem.git
cd Golem

# Build executable
go build -o golem.exe cmd/golem/main.go

# Run Golem
./golem.exe

# Access the web dashboard at # Access the web dashboard at URL_ADDRESS:8080
```

### Configuration
```bash
# Example configuration (config.yaml)
server:
  port: 8080
  static_dir: "web/static"

storage:
  type: "memory" # memory | mysql | postgres
  retention_period: "720h" # 30 days

health_checks:
  default_interval: "60s"
  default_timeout: "10s"
  ```
### API Reference
**Metrics Endpoints:
1. GET /api/metrics: Fetch the latest system metrics.
2. GET /api/metrics/history?duration=1h: Fetch historical system metrics for a specified duration.

**Health Check Endpoints:
1. GET /api/health-checks: List all health checks.
2. POST /api/health-checks: Create a new health check.
3. GET /api/health-checks/{id}: Get details of a specific health check.
4. PUT /api/health-checks/{id}: Update an existing health check.
5. DELETE /api/health-checks/{id}: Delete a health check.
6. GET /api/health-checks/{id}/history?duration=24h: Fetch the history of a health check.

### Health Check Configuration
```bash
{
  "name": "Example Website",
  "type": "http",
  "target": "https://example.com",
  "interval": "60s",
  "timeout": "10s",
  "method": "GET",
  "expectCode": 200
}
```
### Flowchart
image.png