package collector

import (
	"math/rand"
	"runtime"

	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

type (
	MetricsMap map[string]any
)

func CollectMetrics(metricList MetricsMap, count uint64) {
	var (
		counter storage.Counter
		m       runtime.MemStats
	)
	counter = storage.Counter(count)

	runtime.ReadMemStats(&m)

	metricList["Alloc"] = storage.Gauge(m.Alloc)
	metricList["BuckHashSys"] = storage.Gauge(m.BuckHashSys)
	metricList["Frees"] = storage.Gauge(m.Frees)
	metricList["GCCPUFraction"] = storage.Gauge(m.GCCPUFraction)
	metricList["GCSys"] = storage.Gauge(m.GCSys)
	metricList["HeapAlloc"] = storage.Gauge(m.HeapAlloc)
	metricList["HeapIdle"] = storage.Gauge(m.HeapIdle)
	metricList["HeapInuse"] = storage.Gauge(m.HeapInuse)
	metricList["HeapObjects"] = storage.Gauge(m.HeapObjects)
	metricList["HeapReleased"] = storage.Gauge(m.HeapReleased)
	metricList["HeapSys"] = storage.Gauge(m.HeapSys)
	metricList["LastGC"] = storage.Gauge(m.LastGC)
	metricList["Lookups"] = storage.Gauge(m.Lookups)
	metricList["MCacheInuse"] = storage.Gauge(m.MCacheInuse)
	metricList["MCacheSys"] = storage.Gauge(m.MCacheSys)
	metricList["MSpanInuse"] = storage.Gauge(m.MSpanInuse)
	metricList["MSpanSys"] = storage.Gauge(m.MSpanSys)
	metricList["Mallocs"] = storage.Gauge(m.Mallocs)
	metricList["NextGC"] = storage.Gauge(m.NextGC)
	metricList["NumForcedGC"] = storage.Gauge(m.NumForcedGC)
	metricList["NumGC"] = storage.Gauge(m.NumGC)
	metricList["OtherSys"] = storage.Gauge(m.OtherSys)
	metricList["PauseTotalNs"] = storage.Gauge(m.PauseTotalNs)
	metricList["StackInuse"] = storage.Gauge(m.StackInuse)
	metricList["StackSys"] = storage.Gauge(m.StackSys)
	metricList["Sys"] = storage.Gauge(m.Sys)
	metricList["TotalAlloc"] = storage.Gauge(m.TotalAlloc)
	metricList["PollCount"] = counter
	metricList["RandomValue"] = storage.Gauge(rand.Uint64())
}
