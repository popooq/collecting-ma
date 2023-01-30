package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/popooq/collectimg-ma/internal/agent/config"
	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

func main() {
	cfg := config.New()
	hshr := hasher.Mew(cfg.Key)
	sndr := sender.New(hshr)

	var (
		memStat      runtime.MemStats
		counter      int64
		tickerpoll   = time.NewTicker(cfg.PollInterval)
		tickerreport = time.NewTicker(cfg.ReportInterval)
	)

	for {
		select {
		case <-tickerpoll.C:
			runtime.ReadMemStats(&memStat)
			counter++
		case <-tickerreport.C:
			cfg := cfg
			mem := memStat
			counter := storage.Counter(counter)
			random := float64(rand.Uint32())
			sndr.Go(random, "RandomValue", cfg.Address)
			sndr.Go(counter, "PollCount", cfg.Address)
			sndr.Go(float64(mem.Alloc), "Alloc", cfg.Address)
			sndr.Go(float64(mem.BuckHashSys), "BuckHashSys", cfg.Address)
			sndr.Go(float64(mem.Frees), "Frees", cfg.Address)
			sndr.Go(mem.GCCPUFraction, "GCCPUFraction", cfg.Address)
			sndr.Go(float64(mem.GCSys), "GCSys", cfg.Address)
			sndr.Go(float64(mem.HeapAlloc), "HeapAlloc", cfg.Address)
			sndr.Go(float64(mem.HeapIdle), "HeapIdle", cfg.Address)
			sndr.Go(float64(mem.HeapInuse), "HeapInuse", cfg.Address)
			sndr.Go(float64(mem.HeapObjects), "HeapObjects", cfg.Address)
			sndr.Go(float64(mem.HeapReleased), "HeapReleased", cfg.Address)
			sndr.Go(float64(mem.HeapSys), "HeapSys", cfg.Address)
			sndr.Go(float64(mem.LastGC), "LastGC", cfg.Address)
			sndr.Go(float64(mem.Lookups), "Lookups", cfg.Address)
			sndr.Go(float64(mem.MCacheInuse), "MCacheInuse", cfg.Address)
			sndr.Go(float64(mem.MCacheSys), "MCacheSys", cfg.Address)
			sndr.Go(float64(mem.MSpanInuse), "MSpanInuse", cfg.Address)
			sndr.Go(float64(mem.MSpanSys), "MSpanSys", cfg.Address)
			sndr.Go(float64(mem.Mallocs), "Mallocs", cfg.Address)
			sndr.Go(float64(mem.NextGC), "NextGC", cfg.Address)
			sndr.Go(float64(mem.NumForcedGC), "NumForcedGC", cfg.Address)
			sndr.Go(float64(mem.NumGC), "NumGC", cfg.Address)
			sndr.Go(float64(mem.OtherSys), "OtherSys", cfg.Address)
			sndr.Go(float64(mem.PauseTotalNs), "PauseTotalNs", cfg.Address)
			sndr.Go(float64(mem.StackInuse), "StackInuse", cfg.Address)
			sndr.Go(float64(mem.StackSys), "StackSys", cfg.Address)
			sndr.Go(float64(mem.Sys), "Sys", cfg.Address)
			sndr.Go(float64(mem.TotalAlloc), "TotalAlloc", cfg.Address)
		}
	}
}
