package metrics

import (
	"context"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func StartSystemMetricsCollector(ctx context.Context, m *PrometheusMetrics, period time.Duration) {
	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				collectSystemMetrics(m)
			}
		}
	}()
}

func collectSystemMetrics(m *PrometheusMetrics) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		m.SystemCPUUsage.Set(cpuPercent[0])
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	vmStat, err := mem.VirtualMemory()
	if err == nil {
		m.SystemMemoryUsage.Set(float64(vmStat.Used))
	}

	m.ApplicationMemoryUsage.Set(float64(ms.Alloc))
}
