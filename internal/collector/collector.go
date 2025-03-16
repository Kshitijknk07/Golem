package collector

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"

	"Golem/internal/metrics"
	"Golem/internal/storage"
)

type Collector struct {
	storage storage.MetricStorage
}

func NewCollector(storage storage.MetricStorage) *Collector {
	return &Collector{
		storage: storage,
	}
}

func (c *Collector) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics, err := c.collectMetrics()
			if err != nil {
				continue
			}
			c.storage.StoreMetrics(metrics)
		}
	}
}

func (c *Collector) collectMetrics() (metrics.SystemMetrics, error) {
	now := time.Now()

	cpuMetrics, err := c.collectCPUMetrics()
	if err != nil {
		return metrics.SystemMetrics{}, err
	}

	memMetrics, err := c.collectMemoryMetrics()
	if err != nil {
		return metrics.SystemMetrics{}, err
	}

	diskMetrics, err := c.collectDiskMetrics()
	if err != nil {
		return metrics.SystemMetrics{}, err
	}

	networkMetrics, err := c.collectNetworkMetrics()
	if err != nil {
		return metrics.SystemMetrics{}, err
	}

	processMetrics, err := c.collectProcessMetrics()
	if err != nil {
		return metrics.SystemMetrics{}, err
	}

	uptimeMetrics, err := c.collectUptimeMetrics()
	if err != nil {
		return metrics.SystemMetrics{}, err
	}

	return metrics.SystemMetrics{
		Timestamp: now,
		CPU:       cpuMetrics,
		Memory:    memMetrics,
		Disk:      diskMetrics,
		Network:   networkMetrics,
		Process:   processMetrics,
		Uptime:    uptimeMetrics,
	}, nil
}

func (c *Collector) collectCPUMetrics() (metrics.CPUMetrics, error) {
	cpuMetrics := metrics.CPUMetrics{
		PerCoreUsage: make(map[string]float64),
	}

	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return cpuMetrics, err
	}

	if len(percentages) > 0 {
		cpuMetrics.TotalUsage = percentages[0]
	}

	perCPU, err := cpu.Percent(time.Second, true)
	if err != nil {
		return cpuMetrics, err
	}

	for i, p := range perCPU {
		cpuMetrics.PerCoreUsage[string(rune(i))] = p
	}

	loadAvg, err := load.Avg()
	if err == nil {
		cpuMetrics.LoadAverage[0] = loadAvg.Load1
		cpuMetrics.LoadAverage[1] = loadAvg.Load5
		cpuMetrics.LoadAverage[2] = loadAvg.Load15
	}

	return cpuMetrics, nil
}

func (c *Collector) collectMemoryMetrics() (metrics.MemoryMetrics, error) {
	memMetrics := metrics.MemoryMetrics{}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return memMetrics, err
	}

	memMetrics.Total = vmStat.Total
	memMetrics.Used = vmStat.Used
	memMetrics.Free = vmStat.Free
	memMetrics.UsedPercent = vmStat.UsedPercent

	swapStat, err := mem.SwapMemory()
	if err != nil {
		return memMetrics, err
	}

	memMetrics.SwapTotal = swapStat.Total
	memMetrics.SwapUsed = swapStat.Used
	memMetrics.SwapFree = swapStat.Free

	return memMetrics, nil
}

func (c *Collector) collectDiskMetrics() (metrics.DiskMetrics, error) {
	diskMetrics := metrics.DiskMetrics{
		Partitions: []metrics.DiskPartition{},
		IOCounters: make(map[string]metrics.DiskIO),
	}

	partitions, err := disk.Partitions(false)
	if err != nil {
		return diskMetrics, err
	}

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		diskPartition := metrics.DiskPartition{
			Device:      partition.Device,
			Mountpoint:  partition.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		}

		diskMetrics.Partitions = append(diskMetrics.Partitions, diskPartition)
	}

	ioCounters, err := disk.IOCounters()
	if err != nil {
		return diskMetrics, err
	}

	for name, counter := range ioCounters {
		diskMetrics.IOCounters[name] = metrics.DiskIO{
			ReadCount:  counter.ReadCount,
			WriteCount: counter.WriteCount,
			ReadBytes:  counter.ReadBytes,
			WriteBytes: counter.WriteBytes,
			ReadTime:   counter.ReadTime,
			WriteTime:  counter.WriteTime,
		}
	}

	return diskMetrics, nil
}

func (c *Collector) collectNetworkMetrics() (metrics.NetworkMetrics, error) {
	networkMetrics := metrics.NetworkMetrics{
		Interfaces: make(map[string]metrics.NetworkInterface),
	}

	interfaces, err := net.IOCounters(true)
	if err != nil {
		return networkMetrics, err
	}

	for _, netInterface := range interfaces {
		networkMetrics.Interfaces[netInterface.Name] = metrics.NetworkInterface{
			BytesSent:   netInterface.BytesSent,
			BytesRecv:   netInterface.BytesRecv,
			PacketsSent: netInterface.PacketsSent,
			PacketsRecv: netInterface.PacketsRecv,
			Errin:       netInterface.Errin,
			Errout:      netInterface.Errout,
			Dropin:      netInterface.Dropin,
			Dropout:     netInterface.Dropout,
		}
	}

	return networkMetrics, nil
}

func (c *Collector) collectProcessMetrics() ([]metrics.ProcessMetrics, error) {
	var processMetrics []metrics.ProcessMetrics

	processes, err := process.Processes()
	if err != nil {
		return processMetrics, err
	}

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		username, err := p.Username()
		if err != nil {
			username = ""
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}

		memInfo, err := p.MemoryInfo()
		var memUsed uint64 = 0
		if err == nil && memInfo != nil {
			memUsed = memInfo.RSS
		}

		status, err := p.Status()
		if err != nil {
			status = []string{}
		}

		statusStr := ""
		if len(status) > 0 {
			statusStr = status[0]
		}

		createTime, err := p.CreateTime()
		if err != nil {
			createTime = 0
		}

		numThreads, err := p.NumThreads()
		if err != nil {
			numThreads = 0
		}

		ioCounters, err := p.IOCounters()
		processIO := metrics.ProcessIO{}
		if err == nil && ioCounters != nil {
			processIO.ReadCount = ioCounters.ReadCount
			processIO.WriteCount = ioCounters.WriteCount
			processIO.ReadBytes = ioCounters.ReadBytes
			processIO.WriteBytes = ioCounters.WriteBytes
		}

		processMetric := metrics.ProcessMetrics{
			PID:        p.Pid,
			Name:       name,
			Username:   username,
			CPUPercent: cpuPercent,
			MemoryUsed: memUsed,
			Status:     statusStr,
			CreateTime: createTime,
			NumThreads: numThreads,
			IOCounters: processIO,
		}

		processMetrics = append(processMetrics, processMetric)
	}

	return processMetrics, nil
}

func (c *Collector) collectUptimeMetrics() (metrics.UptimeMetrics, error) {
	uptimeMetrics := metrics.UptimeMetrics{}

	hostInfo, err := host.Info()
	if err != nil {
		return uptimeMetrics, err
	}

	uptimeMetrics.Uptime = float64(hostInfo.Uptime)
	uptimeMetrics.BootTime = hostInfo.BootTime

	return uptimeMetrics, nil
}
