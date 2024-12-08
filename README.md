# 📊 Golem Monitoring System  

Welcome to the **Golem Monitoring System**, a production-grade infrastructure monitoring solution! Inspired by industry leaders like Prometheus and Grafana, this system is designed to collect, store, and visualize system and custom metrics. Golem offers a seamless, feature-rich experience for monitoring and managing infrastructure effectively.  

---

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
└── README.md
```
# 🌟 Features

1. Backend API Development  
   - 📈 Metrics Management: Fully functional APIs for fetching, saving, and updating metrics.  
   - 🚨 Alerts Management: Manage custom alert rules and notifications.  
   - 🔐 Secure Authentication: Robust JWT-based authentication and role-based access control (RBAC).  
   - 🔄 Real-Time Updates: WebSocket integration for live updates of metrics and alerts.  

2. Metrics Collection & Data Gathering  
   - ⚙️ System Metrics: Continuous collection of system metrics like CPU and memory usage.  
   - 📊 Prometheus Integration: Collect custom metrics and store time-series data efficiently.  

3. Data Storage & Persistence  
   - 🗃 SQLite Integration: Persistent storage for alerts and configurations.  
   - 🕒 Prometheus Storage: Time-series database for metrics storage and queries.  

4. Alerts & Notification System  
   - 📢 Alert Rules: Pre-configured Prometheus alert rules.  
   - 🔔 Notifications: Integrated with Alertmanager for sending alerts via Slack.  

5. Security & Authentication  
   - 🔒 Secure API Access: JWT-based secure APIs.  
   - 🔐 Data Encryption: Sensitive data is securely encrypted for maximum protection.
     
         Unleash the power of Golem Monitoring System to monitor, manage, and visualize your infrastructure effortlessly. 🎉  

