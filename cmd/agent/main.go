package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/popooq/collectimg-ma/internal/agent/utils/sender"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

type (
	counter int64
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	var (
		m runtime.MemStats
		//	g            gauge
		c            counter
		tickerpoll   = time.NewTicker(pollInterval)
		tickerreport = time.NewTicker(reportInterval)
		// metrics      = collector.MetricsMap{

		// 	"Alloc":         g,
		// 	"BuckHashSys":   g,
		// 	"Frees":         g,
		// 	"GCCPUFraction": g,
		// 	"GCSys":         g,
		// 	"HeapAlloc":     g,
		// 	"HeapIdle":      g,
		// 	"HeapInuse":     g,
		// 	"HeapObjects":   g,
		// 	"HeapReleased":  g,
		// 	"HeapSys":       g,
		// 	"LastGC":        g,
		// 	"Lookups":       g,
		// 	"MCacheInuse":   g,
		// 	"MCacheSys":     g,
		// 	"MSpanInuse":    g,
		// 	"MSpanSys":      g,
		// 	"Mallocs":       g,
		// 	"NextGC":        g,
		// 	"NumForcedGC":   g,
		// 	"NumGC":         g,
		// 	"OtherSys":      g,
		// 	"PauseTotalNs":  g,
		// 	"StackInuse":    g,
		// 	"StackSys":      g,
		// 	"Sys":           g,
		// 	"TotalAlloc":    g,
		// 	"PollCount":     c,
		// 	"RandomValue":   g,
		// }
	)
	for {
		select {
		case <-tickerpoll.C:
			runtime.ReadMemStats(&m)
			//	collector.CollectMetrics(metrics, uint64(c))
			c++
		case <-tickerreport.C:
			//sender.SendMetricsMap(metrics)
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
