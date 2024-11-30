# 📊 Golem Monitoring System

## Project Overview

Welcome to the **Golem Monitoring System**! This project aims to create a robust, production-level monitoring system similar to Prometheus and Grafana. It is designed to collect, store, and visualize system and custom metrics. While the project is functional, it is still far from complete. Several important features and enhancements are yet to be implemented.

## 🚀 Table of Contents

1. [Project Structure](#project-structure)
2. [Current Features](#current-features)
3. [Getting Started](#getting-started)
4. [Running the Project](#running-the-project)
5. [Next Steps](#next-steps)
6. [Contributing](#contributing)
7. [License](#license)

## 🏗 Project Structure

```plaintext
Golem/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── alert.go
│   │   │   ├── metrics.go
│   │   │   ├── middleware.go
│   │   │   └── websockets.go
│   │   ├── db/
│   │   │   └── db.go
│   │   ├── models/
│   │   │   ├── alert.go
│   │   │   └── metric.go
│   │   ├── service/
│   │   │   ├── alert_service.go
│   │   │   └── metrics_service.go
│   │   └── utils/
│   │       └── jwt.go
│   ├── prometheus.yml
│   └── Dockerfile
├── frontend/
│   └── (to be added)
└── README.md
```
## 🌟 Current Features

### Backend

#### API Development

- 📈 Routes for fetching, saving, and updating metrics.
- 🚨 Alerts management.
- 🔐 Authentication and authorization using JWT.
- 🔄 Real-time updates with WebSocket routes.

#### Metrics Collection & Data Gathering

- ⚙️ Integration with Prometheus for system and custom metrics collection.
- 📊 Periodic collection and dynamic updates of CPU and memory usage.

#### Data Storage & Persistence

- 🗃 SQLite database for persistent storage of alerts and configurations.
- 🕒 Prometheus for time-series metrics storage.

#### Alerts & Notification System

- 📢 Configured Prometheus alert rules and Alertmanager for notifications via Slack.

#### Real-Time Dashboard & Frontend

- (Frontend yet to be implemented, planned to use Grafana for visualization.)

#### Security & Authentication

- 🔒 Secure API with JWT and role-based access control (RBAC).
- 🔐 Encryption of sensitive data.
