# Golem - System Monitoring Platform

**Golem** is a complete, production-ready system monitoring solution built in Go. It collects real-time system metrics, performs health checks, and provides a responsive web dashboard for visualization and management.

---

## Features

- **System Metrics**: CPU, memory, disk, network, process, and uptime monitoring.
- **Health Checks**: HTTP, TCP, database, and API endpoint checks with configurable intervals and timeouts.
- **Web Dashboard**: Real-time, interactive dashboard for metrics and health checks.
- **REST API**: Access all metrics and health check data programmatically.
- **Extensible**: Plugin system for custom health checks.
- **In-memory Storage**: Fast, simple storage (data is lost on restart).
- **Easy Setup**: No external dependencies required for basic usage.

---

## Quick Start

### Prerequisites

- Go 1.24 or newer

### Run from Source

```sh
git clone https://github.com/yourusername/Golem.git
cd Golem
go run cmd/golem/main.go
```

Open your browser and go to [http://localhost:8899](http://localhost:8899).

---

## Configuration

- **Port**: Default is `8899`. Change in `.env` or `main.go`.
- **Static Files**: Served from `web/static`.
- **Logging**: Basic info-level logging to stdout.

You can create a `.env` file in `cmd/golem/`:

```
GOLEM_PORT=8899
GOLEM_STATIC_DIR=web/static
GOLEM_LOG_LEVEL=info
```

---

## Usage

- **Web Dashboard**: View system metrics and manage health checks at [http://localhost:8899](http://localhost:8899).
- **API**: Access metrics and health check endpoints under `/api/`.

### Example API Endpoints

- `GET /api/metrics` — Latest system metrics
- `GET /api/metrics/history?duration=1h` — Metrics history
- `GET /api/health-checks` — List health checks
- `POST /api/health-checks` — Create a health check

---

## Project Structure

```
cmd/golem/         # Main entrypoint
internal/api/      # REST API server
internal/collector # Metrics and health check collectors
internal/metrics/  # Data models
internal/storage/  # In-memory storage
web/static/        # Dashboard frontend (HTML/CSS/JS)
```

---

## Extending Golem

- **Plugins**: Implement the `CheckPlugin` interface in Go and register your plugin for custom health checks.
- **Storage**: Only in-memory storage is included. Add your own persistent backend if needed.

---

## Limitations

- **No persistent storage**: All data is lost when the server restarts.
- **No authentication**: The dashboard and API are open by default.
- **Single-node**: No clustering or distributed features.

---

## License

MIT License

---

## Credits

- [gopsutil](https://github.com/shirou/gopsutil) for system metrics
- [gorilla/mux](https://github.com/gorilla/mux) for HTTP routing

---

