package protoreader

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/popooq/collectimg-ma/internal/utils/hasher"
	pb "github.com/popooq/collectimg-ma/proto"
)

type metrics struct {
	value  any
	name   string
	client pb.MetricsClient // пипец временное решение
	hash   hasher.Hash
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
			send(work.client, work.name, work.value, work.hash)
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
func MetricRequest(c pb.MetricsClient, sigs chan os.Signal, poll time.Duration, report time.Duration, rate int, hash hasher.Hash) {
	var (
		memStat runtime.MemStats
		// memoryStat     *mem.VirtualMemoryStat
		// cpuUsage       []float64
		tickerpoll   = time.NewTicker(poll)
		tickerreport = time.NewTicker(report)
		graceperiod  = 15 * time.Second
		shutdown     bool
	)

	shutdown = true

	ctx := context.Background()
	w := newWorker(rate)

	w.start(ctx)
	_, cancel := context.WithTimeout(ctx, graceperiod)
	defer func() {
		w.stop()
		cancel()
	}()
	for shutdown {
		select {
		case <-sigs:
			log.Println("shutdowning the client")
			shutdown = false
		case <-tickerpoll.C:
			runtime.ReadMemStats(&memStat)
			// memoryStat, _ = mem.VirtualMemory()
			// cpuUsage, _ = cpu.Percent(0, false)
		case <-tickerreport.C:
			random := float64(rand.Uint32())
			mem := memStat
			memslice := []metrics{
				{random, "RandomValue", c, hash},
				//{r.pollCount, "PollCount", c, hash},
				{float64(mem.Alloc), "Alloc", c, hash},
				{float64(mem.BuckHashSys), "BuckHashSys", c, hash},
				{float64(mem.Frees), "Frees", c, hash},
				{mem.GCCPUFraction, "GCCPUFraction", c, hash},
				{float64(mem.GCSys), "GCSys", c, hash},
				{float64(mem.HeapAlloc), "HeapAlloc", c, hash},
				{float64(mem.HeapIdle), "HeapIdle", c, hash},
				{float64(mem.HeapInuse), "HeapInuse", c, hash},
				{float64(mem.HeapObjects), "HeapObjects", c, hash},
				{float64(mem.HeapReleased), "HeapReleased", c, hash},
				{float64(mem.HeapSys), "HeapSys", c, hash},
				{float64(mem.LastGC), "LastGC", c, hash},
				{float64(mem.Lookups), "Lookups", c, hash},
				{float64(mem.MCacheInuse), "MCacheInuse", c, hash},
				{float64(mem.MCacheSys), "MCacheSys", c, hash},
				{float64(mem.MSpanInuse), "MSpanInuse", c, hash},
				{float64(mem.MSpanSys), "MSpanSys", c, hash},
				{float64(mem.Mallocs), "Mallocs", c, hash},
				{float64(mem.NextGC), "NextGC", c, hash},
				{float64(mem.NumForcedGC), "NumForcedGC", c, hash},
				{float64(mem.NumGC), "NumGC", c, hash},
				{float64(mem.OtherSys), "OtherSys", c, hash},
				{float64(mem.PauseTotalNs), "PauseTotalNs", c, hash},
				{float64(mem.StackInuse), "StackInuse", c, hash},
				{float64(mem.StackSys), "StackSys", c, hash},
				{float64(mem.Sys), "Sys", c, hash},
				{float64(mem.TotalAlloc), "TotalAlloc", c, hash},
				//	{float64(memoryStat.Total), "TotalMemory", c, hash},
				// {float64(memoryStat.Free), "FreeMemory", c, hash},
				// {float64(cpuUsage[0]), "CPUutilization1", c, hash},
			}
			for _, mem := range memslice {
				w.queueTask(mem)
			}
		}
	}
}

func send(c pb.MetricsClient, name string, value any, hasher hasher.Hash) {

	metric := &pb.Metric{}

	types := strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%T", value), "storage."))

	switch types {
	case "float64":
		assertvalue, ok := value.(float64)
		if !ok {
			log.Printf("conversion failed")
		}
		metric.ID = name
		metric.Value = assertvalue
		metric.Mtype = 0
	case "counter":
		assertvalue, ok := value.(int64)
		if !ok {
			log.Printf("conversion failed")
		}
		metric.ID = name
		metric.Delta = assertvalue
		metric.Mtype = 1
	}

	hash := hasher.HashergRPC(metric)

	err := hasher.HashCheckergRPC(hash, metric)
	if err != nil {
		log.Printf("error: %s", err)
	}

	_, err = c.AddMetric(context.Background(), &pb.AddMetricRequest{
		Metric: metric,
	})
	if err != nil {
		log.Fatalf("err %s", err)
	}
}
