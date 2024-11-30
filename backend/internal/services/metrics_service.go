package service

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var (
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "Current CPU usage percentage",
	})
	memUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mem_usage",
		Help: "Current Memory usage percentage",
	})
)

func InitMetricsCollection() {
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memUsage)

	go collectMetrics()
}

func collectMetrics() {
	for {
		cpuPercents, _ := cpu.Percent(time.Second, false)
		memStats, _ := mem.VirtualMemory()

		if len(cpuPercents) > 0 {
			cpuUsage.Set(cpuPercents[0])
		}
		memUsage.Set(memStats.UsedPercent)

		time.Sleep(10 * time.Second)
	}
}
