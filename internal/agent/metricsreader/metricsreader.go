package metricsreader

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Reader struct {
	sndr         sender.Sender
	tickerpoll   time.Duration
	tickerreport time.Duration
	address      string
	pollCount    storage.Counter
}

func New(sndr sender.Sender, tickerpoll time.Duration, tickerreport time.Duration, address string) *Reader {
	return &Reader{
		sndr:         sndr,
		tickerpoll:   tickerpoll,
		tickerreport: tickerreport,
		address:      address,
	}
}

func (r *Reader) Run() {
	var (
		memStat      runtime.MemStats
		memoryStat   *mem.VirtualMemoryStat
		cpuUsage     []float64
		tickerpoll   = time.NewTicker(r.tickerpoll)
		tickerreport = time.NewTicker(r.tickerreport)
	)

	for {
		select {
		case <-tickerpoll.C:
			runtime.ReadMemStats(&memStat)
			memoryStat, _ = mem.VirtualMemory()
			*&cpuUsage, _ = cpu.Percent(0, false)
			r.pollCount++
		case <-tickerreport.C:
			mem := memStat
			random := float64(rand.Uint32())
			r.sndr.Go(random, "RandomValue", r.address)
			r.sndr.Go(r.pollCount, "PollCount", r.address)
			r.sndr.Go(float64(mem.Alloc), "Alloc", r.address)
			r.sndr.Go(float64(mem.BuckHashSys), "BuckHashSys", r.address)
			r.sndr.Go(float64(mem.Frees), "Frees", r.address)
			r.sndr.Go(mem.GCCPUFraction, "GCCPUFraction", r.address)
			r.sndr.Go(float64(mem.GCSys), "GCSys", r.address)
			r.sndr.Go(float64(mem.HeapAlloc), "HeapAlloc", r.address)
			r.sndr.Go(float64(mem.HeapIdle), "HeapIdle", r.address)
			r.sndr.Go(float64(mem.HeapInuse), "HeapInuse", r.address)
			r.sndr.Go(float64(mem.HeapObjects), "HeapObjects", r.address)
			r.sndr.Go(float64(mem.HeapReleased), "HeapReleased", r.address)
			r.sndr.Go(float64(mem.HeapSys), "HeapSys", r.address)
			r.sndr.Go(float64(mem.LastGC), "LastGC", r.address)
			r.sndr.Go(float64(mem.Lookups), "Lookups", r.address)
			r.sndr.Go(float64(mem.MCacheInuse), "MCacheInuse", r.address)
			r.sndr.Go(float64(mem.MCacheSys), "MCacheSys", r.address)
			r.sndr.Go(float64(mem.MSpanInuse), "MSpanInuse", r.address)
			r.sndr.Go(float64(mem.MSpanSys), "MSpanSys", r.address)
			r.sndr.Go(float64(mem.Mallocs), "Mallocs", r.address)
			r.sndr.Go(float64(mem.NextGC), "NextGC", r.address)
			r.sndr.Go(float64(mem.NumForcedGC), "NumForcedGC", r.address)
			r.sndr.Go(float64(mem.NumGC), "NumGC", r.address)
			r.sndr.Go(float64(mem.OtherSys), "OtherSys", r.address)
			r.sndr.Go(float64(mem.PauseTotalNs), "PauseTotalNs", r.address)
			r.sndr.Go(float64(mem.StackInuse), "StackInuse", r.address)
			r.sndr.Go(float64(mem.StackSys), "StackSys", r.address)
			r.sndr.Go(float64(mem.Sys), "Sys", r.address)
			r.sndr.Go(float64(mem.TotalAlloc), "TotalAlloc", r.address)
			r.sndr.Go(float64(memoryStat.Total), "TotalMemory", r.address)
			r.sndr.Go(float64(memoryStat.Free), "FreeMemory", r.address)
			r.sndr.Go(float64(cpuUsage[0]), "CPUutilization1", r.address)
		}
	}
}
