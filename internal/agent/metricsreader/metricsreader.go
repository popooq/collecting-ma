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

type (
	Reader struct {
		sndr         sender.Sender
		tickerpoll   time.Duration
		tickerreport time.Duration
		address      string
		pollCount    storage.Counter
	}
	metrics struct {
		value any
		name  string
	}
)

func New(sndr sender.Sender, tickerpoll time.Duration, tickerreport time.Duration, address string) *Reader {
	return &Reader{
		sndr:         sndr,
		tickerpoll:   tickerpoll,
		tickerreport: tickerreport,
		address:      address,
	}
}

func (r *Reader) worker(jobs <-chan int, result chan<- int) {
	for j := range jobs {
		r.send()
		result <- j
	}
}

func (r *Reader) Run() {
	const numJobs = 5

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 5; w++ {
		go r.worker(jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}

func (r *Reader) send() {

	var (
		memStat      runtime.MemStats
		memoryStat   *mem.VirtualMemoryStat
		cpuUsage     []float64
		tickerpoll   = time.NewTicker(r.tickerpoll)
		tickerreport = time.NewTicker(r.tickerreport)
		ch           chan sender.Metrics
	)

	for {
		select {
		case <-tickerpoll.C:
			runtime.ReadMemStats(&memStat)
			memoryStat, _ = mem.VirtualMemory()
			cpuUsage, _ = cpu.Percent(0, false)
			r.pollCount++
		case <-tickerreport.C:
			random := float64(rand.Uint32())
			mem := memStat
			memslice := []sender.Metrics{
				{Value: random, Name: "RandomValue"},
				{Value: r.pollCount, Name: "PollCount"},
				{Value: random, Name: "RandomValue"},
				{Value: r.pollCount, Name: "PollCount"},
				{Value: float64(mem.Alloc), Name: "Alloc"},
				{Value: float64(mem.BuckHashSys), Name: "BuckHashSys"},
				{Value: float64(mem.Frees), Name: "Frees"},
				{Value: mem.GCCPUFraction, Name: "GCCPUFraction"},
				{Value: float64(mem.GCSys), Name: "GCSys"},
				{Value: float64(mem.HeapAlloc), Name: "HeapAlloc"},
				{Value: float64(mem.HeapIdle), Name: "HeapIdle"},
				{Value: float64(mem.HeapInuse), Name: "HeapInuse"},
				{Value: float64(mem.HeapObjects), Name: "HeapObjects"},
				{Value: float64(mem.HeapReleased), Name: "HeapReleased"},
				{Value: float64(mem.HeapSys), Name: "HeapSys"},
				{Value: float64(mem.LastGC), Name: "LastGC"},
				{Value: float64(mem.Lookups), Name: "Lookups"},
				{Value: float64(mem.MCacheInuse), Name: "MCacheInuse"},
				{Value: float64(mem.MCacheSys), Name: "MCacheSys"},
				{Value: float64(mem.MSpanInuse), Name: "MSpanInuse"},
				{Value: float64(mem.MSpanSys), Name: "MSpanSys"},
				{Value: float64(mem.Mallocs), Name: "Mallocs"},
				{Value: float64(mem.NextGC), Name: "NextGC"},
				{Value: float64(mem.NumForcedGC), Name: "NumForcedGC"},
				{Value: float64(mem.NumGC), Name: "NumGC"},
				{Value: float64(mem.OtherSys), Name: "OtherSys"},
				{Value: float64(mem.PauseTotalNs), Name: "PauseTotalNs"},
				{Value: float64(mem.StackInuse), Name: "StackInuse"},
				{Value: float64(mem.StackSys), Name: "StackSys"},
				{Value: float64(mem.Sys), Name: "Sys"},
				{Value: float64(mem.TotalAlloc), Name: "TotalAlloc"},
				{Value: float64(memoryStat.Total), Name: "TotalMemory"},
				{Value: float64(memoryStat.Free), Name: "FreeMemory"},
				{Value: float64(cpuUsage[0]), Name: "CPUutilization1"},
			}
			for _, mem := range memslice {
				ch <- mem
			}
			go r.sndr.Go(ch)
		}
	}
}
