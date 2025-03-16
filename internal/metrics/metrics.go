package metrics

import (
	"time"
)

type SystemMetrics struct {
	Timestamp time.Time    `json:"timestamp"`
	CPU       CPUMetrics   `json:"cpu"`
	Memory    MemoryMetrics `json:"memory"`
	Disk      DiskMetrics   `json:"disk"`
	Network   NetworkMetrics `json:"network"`
	Process   []ProcessMetrics `json:"processes"`
	Uptime    UptimeMetrics `json:"uptime"`
}

type CPUMetrics struct {
	TotalUsage    float64            `json:"total_usage"`
	PerCoreUsage  map[string]float64 `json:"per_core_usage"`
	LoadAverage   [3]float64         `json:"load_average"`
}

type MemoryMetrics struct {
	Total       uint64 `json:"total"`
	Used        uint64 `json:"used"`
	Free        uint64 `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	SwapTotal   uint64 `json:"swap_total"`
	SwapUsed    uint64 `json:"swap_used"`
	SwapFree    uint64 `json:"swap_free"`
}

type DiskMetrics struct {
	Partitions []DiskPartition `json:"partitions"`
	IOCounters map[string]DiskIO `json:"io_counters"`
}

type DiskPartition struct {
	Device     string  `json:"device"`
	Mountpoint string  `json:"mountpoint"`
	Total      uint64  `json:"total"`
	Used       uint64  `json:"used"`
	Free       uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskIO struct {
	ReadCount  uint64 `json:"read_count"`
	WriteCount uint64 `json:"write_count"`
	ReadBytes  uint64 `json:"read_bytes"`
	WriteBytes uint64 `json:"write_bytes"`
	ReadTime   uint64 `json:"read_time"`
	WriteTime  uint64 `json:"write_time"`
}

type NetworkMetrics struct {
	Interfaces map[string]NetworkInterface `json:"interfaces"`
}

type NetworkInterface struct {
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
	Errin       uint64 `json:"errin"`
	Errout      uint64 `json:"errout"`
	Dropin      uint64 `json:"dropin"`
	Dropout     uint64 `json:"dropout"`
}

type ProcessMetrics struct {
	PID         int32   `json:"pid"`
	Name        string  `json:"name"`
	Username    string  `json:"username"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryUsed  uint64  `json:"memory_used"`
	Status      string  `json:"status"`
	CreateTime  int64   `json:"create_time"`
	NumThreads  int32   `json:"num_threads"`
	IOCounters  ProcessIO `json:"io_counters,omitempty"`
}

type ProcessIO struct {
	ReadCount  uint64 `json:"read_count"`
	WriteCount uint64 `json:"write_count"`
	ReadBytes  uint64 `json:"read_bytes"`
	WriteBytes uint64 `json:"write_bytes"`
}

type UptimeMetrics struct {
	Uptime  float64 `json:"uptime"`
	BootTime uint64 `json:"boot_time"`
}