package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/popooq/collectimg-ma/internal/agent/utils/sender"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	var (
		m            runtime.MemStats
		c            int64
		tickerpoll   = time.NewTicker(pollInterval)
		tickerreport = time.NewTicker(reportInterval)
	)
	for {
		select {
		case <-tickerpoll.C:
			runtime.ReadMemStats(&m)
			c++
		case <-tickerreport.C:
			sender.SendMetrics(storage.Counter(c), "PollCount")
			sender.SendMetrics(storage.Gauge(m.Alloc), "Alloc")
			sender.SendMetrics(storage.Gauge(m.BuckHashSys), "BuckHashSys")
			sender.SendMetrics(storage.Gauge(m.Frees), "Frees")
			sender.SendMetrics(storage.Gauge(m.GCCPUFraction), "GCCPUFraction")
			sender.SendMetrics(storage.Gauge(m.GCSys), "GCSys")
			sender.SendMetrics(storage.Gauge(m.HeapAlloc), "HeapAlloc")
			sender.SendMetrics(storage.Gauge(m.HeapIdle), "HeapIdle")
			sender.SendMetrics(storage.Gauge(m.HeapInuse), "HeapInuse")
			sender.SendMetrics(storage.Gauge(m.HeapObjects), "HeapObjects")
			sender.SendMetrics(storage.Gauge(m.HeapReleased), "HeapReleased")
			sender.SendMetrics(storage.Gauge(m.HeapSys), "HeapSys")
			sender.SendMetrics(storage.Gauge(m.LastGC), "LastGC")
			sender.SendMetrics(storage.Gauge(m.Lookups), "Lookups")
			sender.SendMetrics(storage.Gauge(m.MCacheInuse), "MCacheInuse")
			sender.SendMetrics(storage.Gauge(m.MCacheSys), "MCacheSys")
			sender.SendMetrics(storage.Gauge(m.MSpanInuse), "MSpanInuse")
			sender.SendMetrics(storage.Gauge(m.MSpanSys), "MSpanSys")
			sender.SendMetrics(storage.Gauge(m.Mallocs), "Mallocs")
			sender.SendMetrics(storage.Gauge(m.NextGC), "NextGC")
			sender.SendMetrics(storage.Gauge(m.NumForcedGC), "NumForcedGC")
			sender.SendMetrics(storage.Gauge(m.NumGC), "NumGC")
			sender.SendMetrics(storage.Gauge(m.OtherSys), "OtherSys")
			sender.SendMetrics(storage.Gauge(m.PauseTotalNs), "PauseTotalNs")
			sender.SendMetrics(storage.Gauge(m.StackInuse), "StackInuse")
			sender.SendMetrics(storage.Gauge(m.StackSys), "StackSys")
			sender.SendMetrics(storage.Gauge(m.Sys), "Sys")
			sender.SendMetrics(storage.Gauge(m.TotalAlloc), "TotalAlloc")
			sender.SendMetrics(storage.Gauge(rand.Uint64()), "RandomValue")
		}
	}
}
