package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/popooq/collectimg-ma/internal/agent/config"
	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/storage"
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
			sender.SendMetrics(storage.Counter(c), "PollCount", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.Alloc), "Alloc", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.BuckHashSys), "BuckHashSys", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.Frees), "Frees", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.GCCPUFraction), "GCCPUFraction", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.GCSys), "GCSys", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.HeapAlloc), "HeapAlloc", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.HeapIdle), "HeapIdle", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.HeapInuse), "HeapInuse", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.HeapObjects), "HeapObjects", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.HeapReleased), "HeapReleased", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.HeapSys), "HeapSys", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.LastGC), "LastGC", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.Lookups), "Lookups", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.MCacheInuse), "MCacheInuse", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.MCacheSys), "MCacheSys", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.MSpanInuse), "MSpanInuse", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.MSpanSys), "MSpanSys", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.Mallocs), "Mallocs", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.NextGC), "NextGC", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.NumForcedGC), "NumForcedGC", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.NumGC), "NumGC", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.OtherSys), "OtherSys", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.PauseTotalNs), "PauseTotalNs", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.StackInuse), "StackInuse", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.StackSys), "StackSys", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.Sys), "Sys", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(m.TotalAlloc), "TotalAlloc", cfg.Address, cfg.Key)
			sender.SendMetrics(float64(rand.Uint64()), "RandomValue", cfg.Address, cfg.Key)
		}
	}
}
