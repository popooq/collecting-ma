package collector

import (
	"math/rand"
	"runtime"
)

type (
	gauge      float64
	counter    uint64
	MetricsMap map[string]any
)

func CollectMetrics(metricList MetricsMap) {
	var (
		c counter
		m runtime.MemStats
	)

	runtime.ReadMemStats(&m)

	metricList["Alloc"] = gauge(m.Alloc)
	metricList["BuckHashSys"] = gauge(m.BuckHashSys)
	metricList["Frees"] = gauge(m.Frees)
	metricList["GCCPUFraction"] = gauge(m.GCCPUFraction)
	metricList["GCSys"] = gauge(m.GCSys)
	metricList["HeapAlloc"] = gauge(m.HeapAlloc)
	metricList["HeapIdle"] = gauge(m.HeapIdle)
	metricList["HeapInuse"] = gauge(m.HeapInuse)
	metricList["HeapObjects"] = gauge(m.HeapObjects)
	metricList["HeapReleased"] = gauge(m.HeapReleased)
	metricList["HeapSys"] = gauge(m.HeapSys)
	metricList["LastGC"] = gauge(m.LastGC)
	metricList["Lookups"] = gauge(m.Lookups)
	metricList["MCacheInuse"] = gauge(m.MCacheInuse)
	metricList["MCacheSys"] = gauge(m.MCacheSys)
	metricList["MSpanInuse"] = gauge(m.MSpanInuse)
	metricList["MSpanSys"] = gauge(m.MSpanSys)
	metricList["Mallocs"] = gauge(m.Mallocs)
	metricList["NextGC"] = gauge(m.NextGC)
	metricList["NumForcedGC"] = gauge(m.NumForcedGC)
	metricList["NumGC"] = gauge(m.NumGC)
	metricList["OtherSys"] = gauge(m.OtherSys)
	metricList["PauseTotalNs"] = gauge(m.PauseTotalNs)
	metricList["StackInuse"] = gauge(m.StackInuse)
	metricList["StackSys"] = gauge(m.StackSys)
	metricList["Sys"] = gauge(m.Sys)
	metricList["TotalAlloc"] = gauge(m.TotalAlloc)
	metricList["PollCount"] = c
	metricList["RandomValue"] = gauge(rand.Uint64())
}
