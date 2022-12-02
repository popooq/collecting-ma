package main

import (
	"time"

	"github.com/popooq/collectimg-ma/internal/agent/utils/collector"
	"github.com/popooq/collectimg-ma/internal/agent/utils/sender"
)

type (
	gauge   float64
	counter uint64
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	var (
		g            gauge
		c            counter
		tickerpoll   = time.NewTicker(pollInterval)
		tickerreport = time.NewTicker(reportInterval)
		metrics      = collector.MetricsMap{
			"Alloc":         g,
			"BuckHashSys":   g,
			"Frees":         g,
			"GCCPUFraction": g,
			"GCSys":         g,
			"HeapAlloc":     g,
			"HeapIdle":      g,
			"HeapInuse":     g,
			"HeapObjects":   g,
			"HeapReleased":  g,
			"HeapSys":       g,
			"LastGC":        g,
			"Lookups":       g,
			"MCacheInuse":   g,
			"MCacheSys":     g,
			"MSpanInuse":    g,
			"MSpanSys":      g,
			"Mallocs":       g,
			"NextGC":        g,
			"NumForcedGC":   g,
			"NumGC":         g,
			"OtherSys":      g,
			"PauseTotalNs":  g,
			"StackInuse":    g,
			"StackSys":      g,
			"Sys":           g,
			"TotalAlloc":    g,
			"PollCount":     c,
			"RandomValue":   g,
		}
	)
	for {
		select {
		case <-tickerpoll.C:
			collector.CollectMetrics(metrics)
			c++
		case <-tickerreport.C:
			sender.SendMetrics(metrics)
		}
	}
}
