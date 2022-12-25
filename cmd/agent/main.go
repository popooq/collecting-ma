package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/popooq/collectimg-ma/internal/agent/config"
	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

func main() {
	cfg := config.NewAgentConfig()
	var (
		m            runtime.MemStats
		c            int64
		tickerpoll   = time.NewTicker(cfg.PollInterval)
		tickerreport = time.NewTicker(cfg.ReportInterval)
	)
	for {
		select {
		case <-tickerpoll.C:
			runtime.ReadMemStats(&m)
			c++
		case <-tickerreport.C:
			sender.SendMetrics(storage.Counter(c), "PollCount", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.Alloc), "Alloc", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.BuckHashSys), "BuckHashSys", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.Frees), "Frees", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.GCCPUFraction), "GCCPUFraction", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.GCSys), "GCSys", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.HeapAlloc), "HeapAlloc", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.HeapIdle), "HeapIdle", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.HeapInuse), "HeapInuse", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.HeapObjects), "HeapObjects", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.HeapReleased), "HeapReleased", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.HeapSys), "HeapSys", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.LastGC), "LastGC", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.Lookups), "Lookups", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.MCacheInuse), "MCacheInuse", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.MCacheSys), "MCacheSys", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.MSpanInuse), "MSpanInuse", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.MSpanSys), "MSpanSys", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.Mallocs), "Mallocs", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.NextGC), "NextGC", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.NumForcedGC), "NumForcedGC", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.NumGC), "NumGC", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.OtherSys), "OtherSys", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.PauseTotalNs), "PauseTotalNs", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.StackInuse), "StackInuse", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.StackSys), "StackSys", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.Sys), "Sys", cfg.Address)
			sender.SendMetrics(storage.Gauge(m.TotalAlloc), "TotalAlloc", cfg.Address)
			sender.SendMetrics(storage.Gauge(rand.Uint64()), "RandomValue", cfg.Address)
		}
	}
}
