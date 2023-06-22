// пакет metricsreader нужен для сбора и отправки метрик из рантайма
package metricsreader

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/storage"
)

// Reader структура хранит информацию о конфигурации и sender
type Reader struct {
	sndr         sender.Sender   // sdnr - отправляет метрики на поле
	tickerpoll   time.Duration   // tickerpoll - частота обновления метрик
	tickerreport time.Duration   // tickerreport - частота отправления метрик
	address      string          // address - адресс сервера куда отправляются метрики
	rate         int             // rate - количество горутин
	pollCount    storage.Counter // pollCount - счетчик отпрпвалений
	shutdown     bool            // shutdown - проверка сисколов
}
type metrics struct {
	value any
	name  string
}
type worker struct {
	workchan   chan metrics
	buffer     int
	wg         *sync.WaitGroup
	cancelFunc context.CancelFunc
	sndr       sender.Sender
}

// WorkerIface
type workerIface interface {
	start(pctx context.Context)
	stop()
	queueTask(mem metrics) error
}

func New(sndr sender.Sender, tickerpoll time.Duration, tickerreport time.Duration, address string, rate int) *Reader {
	return &Reader{
		sndr:         sndr,
		tickerpoll:   tickerpoll,
		tickerreport: tickerreport,
		address:      address,
		rate:         rate,
	}
}

func newWorker(buffer int, sndr sender.Sender) workerIface {
	w := worker{
		workchan: make(chan metrics, buffer),
		buffer:   buffer,
		wg:       new(sync.WaitGroup),
		sndr:     sndr,
	}

	return &w
}

func (w *worker) start(pctx context.Context) {
	ctx, cancelFunc := context.WithCancel(pctx)
	w.cancelFunc = cancelFunc

	for i := 0; i < w.buffer; i++ {
		w.wg.Add(1)
		go w.spawnWorkers(ctx)
	}
}

func (w *worker) stop() {
	close(w.workchan)
	w.cancelFunc()
	w.wg.Wait()
}

func (w *worker) spawnWorkers(ctx context.Context) {
	defer w.wg.Done()

	for work := range w.workchan {
		select {
		case <-ctx.Done():
			return
		default:
			w.sndr.Go(work.value, work.name)
		}
	}
}

func (w *worker) queueTask(mem metrics) error {
	if len(w.workchan) >= w.buffer {
		err := errors.New("workers are busy, try again later")
		return err
	}

	w.workchan <- mem
	log.Println(w.workchan)
	return nil
}

func (r Reader) Run(sigs chan os.Signal) {
	var (
		memStat      runtime.MemStats
		memoryStat   *mem.VirtualMemoryStat
		cpuUsage     []float64
		tickerpoll   = time.NewTicker(r.tickerpoll)
		tickerreport = time.NewTicker(r.tickerreport)
		graceperiod  = 15 * time.Second
	)

	r.shutdown = false

	ctx := context.Background()
	w := newWorker(r.rate, r.sndr)

	w.start(ctx)
	_, cancel := context.WithTimeout(ctx, graceperiod)
	defer func() {
		w.stop()
		cancel()
	}()

	for !r.shutdown {
		select {

		case <-sigs:
			r.shutdown = true

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
