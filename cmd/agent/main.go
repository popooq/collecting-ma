package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type (
	gauge      float64
	counter    uint64
	metricsMap map[string]any
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

var (
	g             gauge
	c             counter
	m             runtime.MemStats
	listOfMetrics metricsMap
	tickerpoll    = time.NewTicker(pollInterval)
	tickerreport  = time.NewTicker(reportInterval)
)

func collectMetrics(metricList metricsMap) metricsMap {
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
	metricList["OtherSys"] = gauge(m.OtherSys)
	metricList["PauseTotalNs"] = gauge(m.PauseTotalNs)
	metricList["StackInuse"] = gauge(m.StackInuse)
	metricList["StackSys"] = gauge(m.StackSys)
	metricList["Sys"] = gauge(m.Sys)
	metricList["TotalAlloc"] = gauge(m.TotalAlloc)
	metricList["PollCount"] = c
	metricList["RandomValue"] = gauge(rand.Uint64())

	return metricList
}

func sendMetrics(data metricsMap) {
	for k, v := range data {

		value := fmt.Sprint(v)
		types := strings.TrimPrefix(fmt.Sprintf("%T", v), "main.")

		endpoint := "http://127.0.0.1:8080/update/" + types + "/" + k + "/" + value + "/"
		resp, err := http.Post(endpoint, "text/plain", nil)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
	}
}
func main() {
	Metrics := metricsMap{
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

	for {
		select {
		case <-tickerpoll.C:
			listOfMetrics = collectMetrics(Metrics)
			c += 1
		case <-tickerreport.C:
			sendMetrics(listOfMetrics)
		}
	}
}
