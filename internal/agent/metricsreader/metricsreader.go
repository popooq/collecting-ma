package metricsreader

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"runtime"
	"sync"
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
		rate         int
		pollCount    storage.Counter
	}
	metrics struct {
		value any
		name  string
	}
	worker struct {
		workchan   chan metrics
		buffer     int
		wg         *sync.WaitGroup
		cancelFunc context.CancelFunc
		reader     Reader
	}
	WorkerIface interface {
		Start(pctx context.Context)
		Stop()
		queueTask(mem metrics) error
	}
)

func New(sndr sender.Sender, tickerpoll time.Duration, tickerreport time.Duration, address string, rate int) *Reader {
	return &Reader{
		sndr:         sndr,
		tickerpoll:   tickerpoll,
		tickerreport: tickerreport,
		address:      address,
		rate:         rate,
	}
}

func NewWorker(buffer int, reader *Reader) WorkerIface {
	w := worker{
		workchan: make(chan metrics, buffer),
		buffer:   buffer,
		wg:       new(sync.WaitGroup),
		reader:   *reader,
	}

	return &w
}

func (w *worker) Start(pctx context.Context) {
	ctx, cancelFunc := context.WithCancel(pctx)
	w.cancelFunc = cancelFunc

	for i := 0; i < w.buffer; i++ {
		w.wg.Add(1)
		log.Printf("cпавн вопркеров")
		go w.spawnWorkers(ctx)
	}
}

func (w *worker) Stop() {
	log.Println("stop workers")
	close(w.workchan)
	w.cancelFunc()
	w.wg.Wait()
	log.Println("all workers exited!")
}

func (w *worker) spawnWorkers(ctx context.Context) {
	defer w.wg.Done()

	for work := range w.workchan {
		log.Printf("work in w.workchan %v", work)
		select {
		case <-ctx.Done():
			log.Printf("ctx done")
			return
		default:
			log.Println("сендер начал работу")
			w.reader.sndr.Go(work.value, work.name)
		}
	}
}

func (w *worker) queueTask(mem metrics) error {
	if len(w.workchan) >= w.buffer {
		log.Println("много воркеров")
		err := errors.New("workers are busy, try again later")
		return err
	}

	log.Printf("mem in qouquueeu %v", mem)
	w.workchan <- mem
	log.Println(w.workchan)
	return nil
}
func (r Reader) Run() {
	var (
		graceperiod  time.Duration = 15 * time.Second
		memStat      runtime.MemStats
		memoryStat   *mem.VirtualMemoryStat
		cpuUsage     []float64
		tickerpoll   = time.NewTicker(r.tickerpoll)
		tickerreport = time.NewTicker(r.tickerreport)
	)

	ctx := context.Background()
	var w worker

	w.Start(ctx)
	_, cancel := context.WithTimeout(ctx, graceperiod)
	defer func() {
		w.Stop()
		cancel()
	}()

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
			memslice := []metrics{
				{random, "RandomValue"},
				{r.pollCount, "PollCount"},
				{float64(mem.Alloc), "Alloc"},
				{float64(mem.BuckHashSys), "BuckHashSys"},
				{float64(mem.Frees), "Frees"},
				{mem.GCCPUFraction, "GCCPUFraction"},
				{float64(mem.GCSys), "GCSys"},
				{float64(mem.HeapAlloc), "HeapAlloc"},
				{float64(mem.HeapIdle), "HeapIdle"},
				{float64(mem.HeapInuse), "HeapInuse"},
				{float64(mem.HeapObjects), "HeapObjects"},
				{float64(mem.HeapReleased), "HeapReleased"},
				{float64(mem.HeapSys), "HeapSys"},
				{float64(mem.LastGC), "LastGC"},
				{float64(mem.Lookups), "Lookups"},
				{float64(mem.MCacheInuse), "MCacheInuse"},
				{float64(mem.MCacheSys), "MCacheSys"},
				{float64(mem.MSpanInuse), "MSpanInuse"},
				{float64(mem.MSpanSys), "MSpanSys"},
				{float64(mem.Mallocs), "Mallocs"},
				{float64(mem.NextGC), "NextGC"},
				{float64(mem.NumForcedGC), "NumForcedGC"},
				{float64(mem.NumGC), "NumGC"},
				{float64(mem.OtherSys), "OtherSys"},
				{float64(mem.PauseTotalNs), "PauseTotalNs"},
				{float64(mem.StackInuse), "StackInuse"},
				{float64(mem.StackSys), "StackSys"},
				{float64(mem.Sys), "Sys"},
				{float64(mem.TotalAlloc), "TotalAlloc"},
				{float64(memoryStat.Total), "TotalMemory"},
				{float64(memoryStat.Free), "FreeMemory"},
				{float64(cpuUsage[0]), "CPUutilization1"},
			}
			for _, mem := range memslice {
				w.queueTask(mem)
			}
		}
	}
}
