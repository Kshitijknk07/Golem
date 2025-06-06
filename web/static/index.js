// Authentication state
let currentUser = null;
let authToken = null;

// DOM Elements
const authSection = document.getElementById("auth-section");
const loginForm = document.getElementById("login-form");
const registerForm = document.getElementById("register-form");
const userManagement = document.getElementById("user-management");
const showRegisterLink = document.getElementById("show-register");
const showLoginLink = document.getElementById("show-login");

// Event Listeners
document.addEventListener("DOMContentLoaded", function () {
  // Check for existing auth token
  const token = localStorage.getItem("authToken");
  if (token) {
    authToken = token;
    validateToken();
  }

  // Form submissions
  document.getElementById("login").addEventListener("submit", handleLogin);
  document
    .getElementById("register")
    .addEventListener("submit", handleRegister);

  // Toggle between login and register forms
  showRegisterLink.addEventListener("click", function (e) {
    e.preventDefault();
    loginForm.style.display = "none";
    registerForm.style.display = "block";
  });

  showLoginLink.addEventListener("click", function (e) {
    e.preventDefault();
    registerForm.style.display = "none";
    loginForm.style.display = "block";
  });

  // Health check form handlers
  document
    .getElementById("add-health-check")
    .addEventListener("click", showAddHealthCheckForm);
  document
    .getElementById("cancel-health-check")
    .addEventListener("click", resetHealthCheckForm);
  document
    .getElementById("save-health-check")
    .addEventListener("click", saveHealthCheck);

  // Initial data load
  if (authToken) {
    loadData();
  }
});

// Authentication Functions
async function validateToken() {
  try {
    const response = await fetch("/api/auth/validate", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    if (response.ok) {
      const user = await response.json();
      setCurrentUser(user);
    } else {
      logout();
    }
  } catch (error) {
    console.error("Error validating token:", error);
    logout();
  }
}

async function handleLogin(e) {
  e.preventDefault();
  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;

  try {
    const response = await fetch("/api/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password }),
    });

    if (response.ok) {
      const data = await response.json();
      authToken = data.token;
      localStorage.setItem("authToken", authToken);
      setCurrentUser(data.user);
      loadData();
    } else {
      const error = await response.json();
      alert(error.message || "Login failed");
    }
  } catch (error) {
    console.error("Login error:", error);
    alert("Login failed: " + error.message);
  }
}

async function handleRegister(e) {
  e.preventDefault();
  const username = document.getElementById("reg-username").value;
  const email = document.getElementById("reg-email").value;
  const password = document.getElementById("reg-password").value;

  try {
    const response = await fetch("/api/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, email, password }),
    });

    if (response.ok) {
      alert("Registration successful! Please login.");
      registerForm.style.display = "none";
      loginForm.style.display = "block";
    } else {
      const error = await response.json();
      alert(error.message || "Registration failed");
    }
  } catch (error) {
    console.error("Registration error:", error);
    alert("Registration failed: " + error.message);
  }
}

function logout() {
  authToken = null;
  currentUser = null;
  localStorage.removeItem("authToken");
  authSection.style.display = "block";
  userManagement.style.display = "none";
  document.querySelectorAll(".card:not(#auth-section)").forEach((card) => {
    card.style.display = "none";
  });
}

function setCurrentUser(user) {
  currentUser = user;
  authSection.style.display = "none";
  document.querySelectorAll(".card:not(#auth-section)").forEach((card) => {
    card.style.display = "block";
  });

  // Show user management only for admin users
  if (user.role === "admin") {
    userManagement.style.display = "block";
    loadUsers();
  } else {
    userManagement.style.display = "none";
  }
}

// User Management Functions
async function loadUsers() {
  if (!currentUser || currentUser.role !== "admin") return;

  try {
    const response = await fetch("/api/users", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    if (response.ok) {
      const users = await response.json();
      const tbody = document.getElementById("users-body");
      tbody.innerHTML = users
        .map(
          (user) => `
                <tr>
                    <td>${user.username}</td>
                    <td>${user.email}</td>
                    <td>${user.role}</td>
                    <td class="user-actions">
                        <button class="btn edit-btn" onclick="editUser('${user.id}')">Edit</button>
                        <button class="btn delete-btn" onclick="deleteUser('${user.id}')">Delete</button>
                    </td>
                </tr>
            `
        )
        .join("");
    }
  } catch (error) {
    console.error("Error loading users:", error);
    alert("Failed to load users: " + error.message);
  }
}

async function editUser(userId) {
  // TODO: Implement user editing
  alert("User editing not implemented yet");
}

async function deleteUser(userId) {
  if (!confirm("Are you sure you want to delete this user?")) return;

  try {
    const response = await fetch(`/api/users/${userId}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    if (response.ok) {
      loadUsers();
    } else {
      const error = await response.json();
      alert(error.message || "Failed to delete user");
    }
  } catch (error) {
    console.error("Error deleting user:", error);
    alert("Failed to delete user: " + error.message);
  }
}

// Data Loading Functions
async function loadData() {
  if (!authToken) return;

  try {
    await Promise.all([loadMetrics(), loadProcesses(), loadHealthChecks()]);
  } catch (error) {
    console.error("Error loading data:", error);
    alert("Failed to load data: " + error.message);
  }
}

async function loadMetrics() {
  try {
    const response = await fetch("/api/metrics/latest", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    if (response.ok) {
      const metrics = await response.json();
      updateMetricsDisplay(metrics);
    }
  } catch (error) {
    console.error("Error loading metrics:", error);
  }
}

async function loadProcesses() {
  try {
    const response = await fetch("/api/metrics/processes", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    if (response.ok) {
      const processes = await response.json();
      updateProcessesDisplay(processes);
    }
  } catch (error) {
    console.error("Error loading processes:", error);
  }
}

async function loadHealthChecks() {
  try {
    const response = await fetch("/api/health-checks", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    if (response.ok) {
      const checks = await response.json();
      updateHealthChecksDisplay(checks);
    }
  } catch (error) {
    console.error("Error loading health checks:", error);
  }
}

// Display Update Functions
function updateMetricsDisplay(metrics) {
  document.getElementById("cpu-usage").textContent = `${metrics.cpu.usage}%`;
  document.getElementById(
    "memory-usage"
  ).textContent = `${metrics.memory.used_percent}%`;
  document.getElementById(
    "disk-usage"
  ).textContent = `${metrics.disk.used_percent}%`;
  document.getElementById("uptime").textContent = formatUptime(metrics.uptime);

  // Update detailed sections
  updateCPUDetails(metrics.cpu);
  updateMemoryDetails(metrics.memory);
  updateDiskDetails(metrics.disk);
  updateNetworkDetails(metrics.network);
}

function updateProcessesDisplay(processes) {
  const tbody = document.getElementById("processes-body");
  tbody.innerHTML = processes
    .map(
      (process) => `
        <tr>
            <td>${process.pid}</td>
            <td>${process.name}</td>
            <td>${process.cpu_percent.toFixed(1)}%</td>
            <td>${formatBytes(process.memory_usage)}</td>
            <td>${process.status}</td>
        </tr>
    `
    )
    .join("");
}

function updateHealthChecksDisplay(checks) {
  const tbody = document.getElementById("health-checks-body");
  tbody.innerHTML = checks
    .map(
      (check) => `
        <tr>
            <td>${check.name}</td>
            <td>${check.type}</td>
            <td>${check.target}</td>
            <td class="status-${check.status.toLowerCase()}">${
        check.status
      }</td>
            <td>${check.response_time}ms</td>
            <td>${new Date(check.last_check).toLocaleString()}</td>
            <td>
                <button class="btn" onclick="editHealthCheck('${
                  check.id
                }')">Edit</button>
                <button class="btn delete-btn" onclick="deleteHealthCheck('${
                  check.id
                }')">Delete</button>
            </td>
        </tr>
    `
    )
    .join("");

  // Update summary metrics
  const total = checks.length;
  const healthy = checks.filter((c) => c.status === "HEALTHY").length;
  const warning = checks.filter((c) => c.status === "WARNING").length;
  const down = checks.filter((c) => c.status === "DOWN").length;

  document.getElementById("total-services").textContent = total;
  document.getElementById("healthy-services").textContent = healthy;
  document.getElementById("warning-services").textContent = warning;
  document.getElementById("down-services").textContent = down;
}

// Utility Functions
function formatBytes(bytes) {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}

function formatUptime(seconds) {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  return `${days}d ${hours}h ${minutes}m`;
}

// Health Check Functions
function showAddHealthCheckForm() {
  document.getElementById("health-check-form").style.display = "block";
  document.getElementById("add-health-check").style.display = "none";
}

function resetHealthCheckForm() {
  document.getElementById("health-check-form").style.display = "none";
  document.getElementById("add-health-check").style.display = "block";
  document.getElementById("check-form").reset();
}

async function saveHealthCheck() {
  const name = document.getElementById("check-name").value;
  const type = document.getElementById("check-type").value;
  const target = document.getElementById("check-url").value;
  const interval =
    parseInt(document.getElementById("check-interval").value) * 1000000000;
  const timeout =
    parseInt(document.getElementById("check-timeout").value) * 1000000000;

  if (!name || !target) {
    alert("Name and URL/Host are required");
    return;
  }

  const healthCheck = {
    name,
    type,
    target,
    interval,
    timeout,
  };

  try {
    const response = await fetch("/api/health-checks", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
      body: JSON.stringify(healthCheck),
    });

    if (response.ok) {
      resetHealthCheckForm();
      loadHealthChecks();
    } else {
      const error = await response.json();
      alert(error.message || "Failed to create health check");
    }
  } catch (error) {
    console.error("Error creating health check:", error);
    alert("Failed to create health check: " + error.message);
  }
}

async function deleteHealthCheck(id) {
  if (!confirm("Are you sure you want to delete this health check?")) return;

  try {
    const response = await fetch(`/api/health-checks/${id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    if (response.ok) {
      loadHealthChecks();
    } else {
      const error = await response.json();
      alert(error.message || "Failed to delete health check");
    }
  } catch (error) {
    console.error("Error deleting health check:", error);
    alert("Failed to delete health check: " + error.message);
  }
}

// Start periodic updates
setInterval(() => {
  if (authToken) {
    loadData();
  }
}, 30000);
