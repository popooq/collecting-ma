package protoreader

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"

	pb "github.com/popooq/collectimg-ma/proto"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type metrics struct {
	value  any
	name   string
	client pb.MetricsClient // пипец временное решение
}

type worker struct {
	workchan   chan metrics
	buffer     int
	wg         *sync.WaitGroup
	cancelFunc context.CancelFunc
}

// WorkerIface
type workerIface interface {
	start(pctx context.Context)
	stop()
	queueTask(mem metrics) error
}

func newWorker(buffer int) workerIface {
	w := worker{
		workchan: make(chan metrics, buffer),
		buffer:   buffer,
		wg:       new(sync.WaitGroup),
	}

	return &w
}

func (w *worker) start(pctx context.Context) {
	ctx, cancelfunc := context.WithCancel(pctx)
	w.cancelFunc = cancelfunc

	for i := 0; i <= w.buffer; i++ {
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
			send(work.client, work.name, work.value)
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
func MetricRequest(c pb.MetricsClient) {
	var (
		memStat        runtime.MemStats
		memoryStat     *mem.VirtualMemoryStat
		cpuUsage       []float64
		tickerpoll     = time.NewTicker(5)
		tickerreport   = time.NewTicker(10)
		graceperiod    = 15 * time.Second
		goroutinestest = 6
	)

	ctx := context.Background()
	w := newWorker(goroutinestest)

	w.start(ctx)
	_, cancel := context.WithTimeout(ctx, graceperiod)
	defer func() {
		w.stop()
		cancel()
	}()
	for {
		select {
		case <-tickerpoll.C:
			runtime.ReadMemStats(&memStat)
			memoryStat, _ = mem.VirtualMemory()
			cpuUsage, _ = cpu.Percent(0, false)
		case <-tickerreport.C:
			random := float64(rand.Uint32())
			mem := memStat
			memslice := []metrics{
				{random, "RandomValue", c},
				//{r.pollCount, "PollCount", c},
				{float64(mem.Alloc), "Alloc", c},
				{float64(mem.BuckHashSys), "BuckHashSys", c},
				{float64(mem.Frees), "Frees", c},
				{mem.GCCPUFraction, "GCCPUFraction", c},
				{float64(mem.GCSys), "GCSys", c},
				{float64(mem.HeapAlloc), "HeapAlloc", c},
				{float64(mem.HeapIdle), "HeapIdle", c},
				{float64(mem.HeapInuse), "HeapInuse", c},
				{float64(mem.HeapObjects), "HeapObjects", c},
				{float64(mem.HeapReleased), "HeapReleased", c},
				{float64(mem.HeapSys), "HeapSys", c},
				{float64(mem.LastGC), "LastGC", c},
				{float64(mem.Lookups), "Lookups", c},
				{float64(mem.MCacheInuse), "MCacheInuse", c},
				{float64(mem.MCacheSys), "MCacheSys", c},
				{float64(mem.MSpanInuse), "MSpanInuse", c},
				{float64(mem.MSpanSys), "MSpanSys", c},
				{float64(mem.Mallocs), "Mallocs", c},
				{float64(mem.NextGC), "NextGC", c},
				{float64(mem.NumForcedGC), "NumForcedGC", c},
				{float64(mem.NumGC), "NumGC", c},
				{float64(mem.OtherSys), "OtherSys", c},
				{float64(mem.PauseTotalNs), "PauseTotalNs", c},
				{float64(mem.StackInuse), "StackInuse", c},
				{float64(mem.StackSys), "StackSys", c},
				{float64(mem.Sys), "Sys", c},
				{float64(mem.TotalAlloc), "TotalAlloc", c},
				//	{float64(memoryStat.Total), "TotalMemory", c},
				{float64(memoryStat.Free), "FreeMemory", c},
				{float64(cpuUsage[0]), "CPUutilization1", c},
			}
			for _, mem := range memslice {
				w.queueTask(mem)
			}
		}
	}
}

func send(c pb.MetricsClient, name string, value any) {

	metric := &pb.Metric{}

	types := strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%T", value), "storage."))

	metric.Mtype = types

	switch metric.Mtype {
	case "float64":
		assertvalue, ok := value.(float64)
		if !ok {
			log.Printf("conversion failed")
		}
		metric.ID = name
		metric.Value = assertvalue
	case "counter":
		assertvalue, ok := value.(int64)
		if !ok {
			log.Printf("conversion failed")
		}
		metric.ID = name
		metric.Delta = assertvalue
	}

	resp, err := c.AddMetric(context.Background(), &pb.AddMetricRequest{
		Metric: metric,
	})
	if err != nil {
		log.Fatalf("err %s", err)
	}
	if resp.Error != "" {
		log.Fatalf("resp err %s", resp.Error)
	}
}
