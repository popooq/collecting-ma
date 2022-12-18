package collector

import (
	"math/rand"
	"runtime"
)

type (
	Gauge      float64
	Counter    uint64
	MetricsMap map[string]any
)

func CollectMetrics(metricList MetricsMap, count uint64) {
	var (
		c Counter
		m runtime.MemStats
	)
	c = Counter(count)

	runtime.ReadMemStats(&m)

	metricList["Alloc"] = Gauge(m.Alloc)
	metricList["BuckHashSys"] = Gauge(m.BuckHashSys)
	metricList["Frees"] = Gauge(m.Frees)
	metricList["GCCPUFraction"] = Gauge(m.GCCPUFraction)
	metricList["GCSys"] = Gauge(m.GCSys)
	metricList["HeapAlloc"] = Gauge(m.HeapAlloc)
	metricList["HeapIdle"] = Gauge(m.HeapIdle)
	metricList["HeapInuse"] = Gauge(m.HeapInuse)
	metricList["HeapObjects"] = Gauge(m.HeapObjects)
	metricList["HeapReleased"] = Gauge(m.HeapReleased)
	metricList["HeapSys"] = Gauge(m.HeapSys)
	metricList["LastGC"] = Gauge(m.LastGC)
	metricList["Lookups"] = Gauge(m.Lookups)
	metricList["MCacheInuse"] = Gauge(m.MCacheInuse)
	metricList["MCacheSys"] = Gauge(m.MCacheSys)
	metricList["MSpanInuse"] = Gauge(m.MSpanInuse)
	metricList["MSpanSys"] = Gauge(m.MSpanSys)
	metricList["Mallocs"] = Gauge(m.Mallocs)
	metricList["NextGC"] = Gauge(m.NextGC)
	metricList["NumForcedGC"] = Gauge(m.NumForcedGC)
	metricList["NumGC"] = Gauge(m.NumGC)
	metricList["OtherSys"] = Gauge(m.OtherSys)
	metricList["PauseTotalNs"] = Gauge(m.PauseTotalNs)
	metricList["StackInuse"] = Gauge(m.StackInuse)
	metricList["StackSys"] = Gauge(m.StackSys)
	metricList["Sys"] = Gauge(m.Sys)
	metricList["TotalAlloc"] = Gauge(m.TotalAlloc)
	metricList["PollCount"] = c
	metricList["RandomValue"] = Gauge(rand.Uint64())
}
