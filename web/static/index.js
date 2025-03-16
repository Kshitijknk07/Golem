document.addEventListener('DOMContentLoaded', function () {
    document.getElementById('add-health-check').addEventListener('click', showAddHealthCheckForm);
    document.getElementById('cancel-health-check').addEventListener('click', resetHealthCheckForm);
    document.getElementById('save-health-check').addEventListener('click', saveHealthCheck);

    updateMetrics();
    updateHealthChecks();

    setInterval(updateMetrics, 5000);
    setInterval(updateHealthChecks, 30000);
});

function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

function formatUptime(seconds) {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${days}d ${hours}h ${minutes}m`;
}

function updateMetrics() {
    fetch('/api/metrics')
        .then(response => response.json())
        .then(data => {
            // System Overview updates
            document.getElementById('cpu-usage').textContent = `${data.cpu.total_usage.toFixed(1)}%`;
            document.getElementById('memory-usage').textContent = `${data.memory.used_percent.toFixed(1)}%`;

            // Disk usage calculation
            let avgDiskUsage = 0;
            if (data.disk.partitions.length > 0) {
                avgDiskUsage = data.disk.partitions.reduce((sum, p) => sum + p.used_percent, 0) / data.disk.partitions.length;
            }
            document.getElementById('disk-usage').textContent = `${avgDiskUsage.toFixed(1)}%`;
            document.getElementById('uptime').textContent = formatUptime(data.uptime.uptime);

            // Detailed metrics updates
            updateCpuDetails(data.cpu);
            updateMemoryDetails(data.memory);
            updateDiskDetails(data.disk);
            updateNetworkDetails(data.network);
            updateProcessDetails(data.processes);
        })
        .catch(console.error);
}

function updateCpuDetails(cpuData) {
    let html = `<p>Load Average: ${cpuData.load_average.map(v => v.toFixed(2)).join(', ')}</p><div class="metric-row">`;
    Object.entries(cpuData.per_core_usage).forEach(([core, usage]) => {
        html += `
            <div class="metric-box">
                <div class="metric-title">Core ${core}</div>
                <div class="metric-value">${usage.toFixed(1)}%</div>
            </div>`;
    });
    document.getElementById('cpu-details').innerHTML = html + '</div>';
}

function updateMemoryDetails(memoryData) {
    const html = `
        <div class="metric-row">
            <div class="metric-box">
                <div class="metric-title">Total Memory</div>
                <div class="metric-value">${formatBytes(memoryData.total)}</div>
            </div>
            <div class="metric-box">
                <div class="metric-title">Used Memory</div>
                <div class="metric-value">${formatBytes(memoryData.used)}</div>
            </div>
            <div class="metric-box">
                <div class="metric-title">Free Memory</div>
                <div class="metric-value">${formatBytes(memoryData.free)}</div>
            </div>
            <div class="metric-box">
                <div class="metric-title">Swap Used</div>
                <div class="metric-value">${formatBytes(memoryData.swap_used)}</div>
            </div>
        </div>`;
    document.getElementById('memory-details').innerHTML = html;
}

function updateDiskDetails(diskData) {
    let html = '<div class="metric-row">';
    diskData.partitions.forEach(partition => {
        html += `
            <div class="metric-box">
                <div class="metric-title">${partition.mountpoint}</div>
                <div class="metric-value">${partition.used_percent.toFixed(1)}%</div>
                <div>${formatBytes(partition.used)} / ${formatBytes(partition.total)}</div>
            </div>`;
    });
    document.getElementById('disk-details').innerHTML = html + '</div>';
}

function updateNetworkDetails(networkData) {
    let html = '<div class="metric-row">';
    Object.entries(networkData.interfaces).forEach(([name, iface]) => {
        html += `
            <div class="metric-box">
                <div class="metric-title">${name}</div>
                <div>Sent: ${formatBytes(iface.bytes_sent)}</div>
                <div>Received: ${formatBytes(iface.bytes_recv)}</div>
                <div>Packets Sent: ${iface.packets_sent}</div>
                <div>Packets Received: ${iface.packets_recv}</div>
            </div>`;
    });
    document.getElementById('network-details').innerHTML = html + '</div>';
}

function updateProcessDetails(processes) {
    let html = '';
    processes.sort((a, b) => b.cpu_percent - a.cpu_percent)
        .slice(0, 10)
        .forEach(proc => {
            html += `
                <tr>
                    <td>${proc.pid}</td>
                    <td>${proc.name}</td>
                    <td>${proc.cpu_percent.toFixed(1)}%</td>
                    <td>${formatBytes(proc.memory_used)}</td>
                    <td>${proc.status}</td>
                </tr>`;
        });
    document.getElementById('processes-body').innerHTML = html || '<tr><td colspan="5">No process data available</td></tr>';
}

// Health Checks Management
function showAddHealthCheckForm() {
    document.getElementById('health-check-form').style.display = 'block';
    document.getElementById('add-health-check').style.display = 'none';
}

function hideAddHealthCheckForm() {
    document.getElementById('health-check-form').style.display = 'none';
    document.getElementById('add-health-check').style.display = 'block';
}

function resetHealthCheckForm() {
    hideAddHealthCheckForm();
    document.getElementById('check-form').reset();
    const saveButton = document.getElementById('save-health-check');
    saveButton.textContent = 'Save';
    saveButton.onclick = saveHealthCheck;
}

function saveHealthCheck() {
    const name = document.getElementById('check-name').value;
    const type = document.getElementById('check-type').value;
    const target = document.getElementById('check-url').value;
    const interval = parseInt(document.getElementById('check-interval').value) * 1e9;
    const timeout = parseInt(document.getElementById('check-timeout').value) * 1e9;

    if (!name || !target) {
        alert('Name and URL/Host are required');
        return;
    }

    const healthCheck = {
        name: name,
        type: type,
        target: target,
        interval: interval,
        timeout: timeout,
        enabled: true
    };

    fetch('/api/health-checks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(healthCheck)
    })
        .then(response => {
            if (!response.ok) throw new Error('Failed to save health check');
            resetHealthCheckForm();
            updateHealthChecks();
        })
        .catch(error => {
            console.error('Error saving health check:', error);
            alert(`Error saving health check: ${error.message}`);
        });
}

function updateHealthChecks() {
    fetch('/api/health-checks')
        .then(response => response.json())
        .then(data => {
            updateHealthChecksSummary(data);
            updateHealthChecksTable(data.checks);
        })
        .catch(console.error);
}

function updateHealthChecksSummary(data) {
    document.getElementById('total-services').textContent = data.checks.length;
    document.getElementById('healthy-services').textContent = data.checks.filter(c => c.status === 'up').length;
    document.getElementById('warning-services').textContent = data.checks.filter(c => c.status === 'warning').length;
    document.getElementById('down-services').textContent = data.checks.filter(c => c.status === 'down').length;
}

function updateHealthChecksTable(checks) {
    let html = '';
    checks.forEach(check => {
        html += `
            <tr>
                <td>${check.name}</td>
                <td>${check.type}</td>
                <td>${check.target}</td>
                <td class="status-${check.status}">${check.status.toUpperCase()}</td>
                <td>${(check.responseTime / 1e6).toFixed(2)}ms</td>
                <td>${new Date(check.lastChecked).toLocaleString()}</td>
                <td>
                    <button onclick="editHealthCheck('${check.id}')" class="btn">Edit</button>
                    <button onclick="deleteHealthCheck('${check.id}')" class="btn">Delete</button>
                </td>
            </tr>`;
    });
    document.getElementById('health-checks-body').innerHTML = html || '<tr><td colspan="7">No health checks configured</td></tr>';
}

function editHealthCheck(id) {
    fetch(`/api/health-checks/${id}`)
        .then(response => response.json())
        .then(check => {
            document.getElementById('check-name').value = check.name;
            document.getElementById('check-type').value = check.type;
            document.getElementById('check-url').value = check.target;
            document.getElementById('check-interval').value = Math.floor(check.interval / 1e9);
            document.getElementById('check-timeout').value = Math.floor(check.timeout / 1e9);

            showAddHealthCheckForm();
            const saveButton = document.getElementById('save-health-check');
            saveButton.textContent = 'Update';
            saveButton.onclick = () => updateExistingHealthCheck(id);
        })
        .catch(error => {
            console.error('Error fetching health check:', error);
            alert(`Error fetching health check: ${error.message}`);
        });
}

function updateExistingHealthCheck(id) {
    const name = document.getElementById('check-name').value;
    const type = document.getElementById('check-type').value;
    const target = document.getElementById('check-url').value;
    const interval = parseInt(document.getElementById('check-interval').value) * 1e9;
    const timeout = parseInt(document.getElementById('check-timeout').value) * 1e9;

    if (!name || !target) {
        alert('Name and URL/Host are required');
        return;
    }

    const healthCheck = {
        name: name,
        type: type,
        target: target,
        interval: interval,
        timeout: timeout,
        enabled: true
    };

    fetch(`/api/health-checks/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(healthCheck)
    })
        .then(response => {
            if (!response.ok) throw new Error('Failed to update health check');
            resetHealthCheckForm();
            updateHealthChecks();
        })
        .catch(error => {
            console.error('Error updating health check:', error);
            alert(`Error updating health check: ${error.message}`);
        });
}

function deleteHealthCheck(id) {
    if (confirm('Are you sure you want to delete this health check?')) {
        fetch(`/api/health-checks/${id}`, { method: 'DELETE' })
            .then(response => {
                if (!response.ok) throw new Error('Failed to delete health check');
                updateHealthChecks();
            })
            .catch(error => {
                console.error('Error deleting health check:', error);
                alert(`Error deleting health check: ${error.message}`);
            });
    }
}

