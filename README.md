# Golem - Distributed System Monitoring Platform

![Go Version](https://img.shields.io/badge/go-1.24%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)

Golem is a monitoring solution for distributed systems, providing real-time metrics collection, health checking, and visualization capabilities.

## Features

### 📊 System Metrics Monitoring
- CPU, Memory, and Disk Usage
- Network I/O Statistics
- Process Monitoring
- Historical Data Storage

### 🩺 Health Check System
- HTTP/HTTPS Endpoint Checks
- TCP Port Availability
- Database Connectivity (MySQL, PostgreSQL, SQLite)
- Customizable Check Intervals

### 🌐 Web Dashboard
- Real-time Metrics Visualization
- Health Check Status Overview
- Historical Trend Analysis

### 🔌 Extensible Architecture
- Plugin System for Custom Checks
- Multiple Storage Backends
- REST API for Integration

---

## Getting Started

### Prerequisites
- Go 1.24+
- MySQL/PostgreSQL (optional for persistent storage)

### Installation
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

---

## Configuration

Create a `config.yaml` file for Golem:
```yaml
server:
  port: 8080
  static_dir: "web/static"

storage:
  type: "memory"  # Options: memory | mysql | postgres
  retention_period: "720h"  # 30 days

health_checks:
  default_interval: "60s"
  default_timeout: "10s"
```

---

## API Reference

### **Metrics Endpoints**
- **GET** `/api/metrics` - Fetch the latest system metrics.
- **GET** `/api/metrics/history?duration=1h` - Fetch historical system metrics for a specified duration.

### **Health Check Endpoints**
- **GET** `/api/health-checks` - List all health checks.
- **POST** `/api/health-checks` - Create a new health check.
- **GET** `/api/health-checks/{id}` - Get details of a specific health check.
- **PUT** `/api/health-checks/{id}` - Update an existing health check.
- **DELETE** `/api/health-checks/{id}` - Delete a health check.
- **GET** `/api/health-checks/{id}/history?duration=24h` - Fetch the history of a health check.

---

## Health Check Configuration Example

Example JSON payload for creating a new health check:
```json
{
  "name": "Example Website",
  "type": "http",
  "target": "https://example.com",
  "interval": "60s",
  "timeout": "10s",
  "method": "GET",
  "expectCode": 200
}
