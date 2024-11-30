package service

import (
	"log"
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
		cpuPercents, err := cpu.Percent(time.Second, false)
		if err != nil {
			log.Printf("Error fetching CPU usage: %v", err)
		} else if len(cpuPercents) > 0 {
			cpuUsage.Set(cpuPercents[0])
			log.Printf("CPU Usage: %.2f%%", cpuPercents[0])
		}

		memStats, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("Error fetching memory usage: %v", err)
		} else {
			memUsage.Set(memStats.UsedPercent)
			log.Printf("Memory Usage: %.2f%%", memStats.UsedPercent)
		}

		time.Sleep(10 * time.Second)
	}
}
