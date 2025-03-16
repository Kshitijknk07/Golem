# Golem Monitoring System

A lightweight system monitoring tool built with Go, providing real-time metrics collection and visualization for your computer's performance.

## Overview

Golem is a simple yet effective monitoring system that provides:
- Real-time system metrics collection (CPU, memory, disk, network)
- Process monitoring with detailed statistics
- Clean and intuitive web dashboard
- RESTful API for metrics access

## Architecture

```
backend/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers and routes
│   ├── db/             # Database operations
│   ├── models/         # Data structures
│   ├── services/       # Business logic
│   └── utils/          # Utility functions
└── prometheus.yml      # Prometheus configuration
```

## Features

### Metrics Management
- RESTful API for CRUD operations on metrics
- Real-time metric streaming via WebSocket
- Persistent storage with SQLite
- Prometheus integration for metrics collection

### Alert System
- Custom alert creation and management
- Severity-based alert categorization
- Real-time alert notifications
- Persistent alert history

### Security
- JWT-based authentication
- Secure WebSocket connections
- Environment-based configuration
- Input validation and sanitization

## API Endpoints

### Metrics
- `GET /metrics` - Fetch all metrics
- `POST /metrics` - Create new metric
- `PUT /metrics` - Update existing metric
- `GET /ws/metrics` - WebSocket stream for real-time updates

### Alerts
- `GET /alerts` - Fetch all alerts
- `POST /alerts` - Create new alert

### Authentication
- `POST /login` - Authenticate and receive JWT token

## Getting Started

### Prerequisites
- Go 1.18 or higher
- SQLite3

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/golem.git
cd golem/backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables (optional):
```bash
export JWT_SECRET_KEY=your_secure_key_here
```

4. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on port 4000.

## Development

### Project Structure
- `internal/api`: HTTP handlers and route definitions
- `internal/db`: Database initialization and operations
- `internal/models`: Data structures and database models
- `internal/services`: Business logic implementation
- `internal/utils`: Utility functions and helpers

### Database Schema

#### Metrics Table
```sql
CREATE TABLE metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    value REAL NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Alerts Table
```sql
CREATE TABLE alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message TEXT NOT NULL,
    severity TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Security

- All API endpoints are protected with JWT authentication
- WebSocket connections are secured
- Database operations are properly sanitized
- Environment variables for sensitive configuration
