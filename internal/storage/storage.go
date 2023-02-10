package storage

import (
	"fmt"
	"sync"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

type (
	Storage interface {
		InsertMetric(name string, value float64)
		CountCounterMetric(name string, value int64)
		GetMetricGauge(name string) (float64, error)
		GetAllMetrics() []string
		GetMetricCounter(name string) (int64, error)
		GetMetricJSONGauge(name string) (*float64, error)
		GetMetricJSONCounter(name string) (*int64, error)
		InsertMetrics(metric encoder.Encode) error
	}

	Keeper interface {
		SaveMetric(metric *encoder.Encode) error
		SaveAllMetrics(metric encoder.Encode) error
		LoadMetrics() ([]encoder.Encode, error)
		KeeperCheck() error
	}

	MetricsStorage struct {
		Keeper         Keeper
		cfg            config.Config
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
		mu             sync.Mutex
	}
	Counter int64
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

func New(Keeper Keeper, cfg config.Config) *MetricsStorage {
	return &MetricsStorage{
		mu:             sync.Mutex{},
		MetricsGauge:   make(map[string]float64),
		MetricsCounter: make(map[string]int64),
		Keeper:         Keeper,
		cfg:            cfg,
	}

}

func (ms *MetricsStorage) InsertMetric(name string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.MetricsGauge[name] = value
}

func (ms *MetricsStorage) CountCounterMetric(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.MetricsCounter[name] += value
}

func (ms *MetricsStorage) GetMetricGauge(name string) (float64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	value, ok := ms.MetricsGauge[name]
	if !ok {
		err := fmt.Errorf("metric %s doesn't exist", name)

		return 0, err
	}

	return value, nil
}

func (ms *MetricsStorage) GetMetricJSONGauge(name string) (*float64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	value, ok := ms.MetricsGauge[name]
	if !ok {
		err := fmt.Errorf("metric %s doesn't exist", name)

		return nil, err
	}

	return &value, nil
}

func (ms *MetricsStorage) GetMetricCounter(name string) (int64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	value, ok := ms.MetricsCounter[name]

	if !ok {
		err := fmt.Errorf("metric %s doesn't exist", name)

		return 0, err
	}

	return value, nil
}

func (ms *MetricsStorage) GetMetricJSONCounter(name string) (*int64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	uvalue, ok := ms.MetricsCounter[name]
	if !ok {
		err := fmt.Errorf("metric %s doesn't exist", name)

		return nil, err
	}

	value := uvalue

	return &value, nil
}

func (ms *MetricsStorage) InsertMetrics(metric encoder.Encode) error {
	switch {
	case metric.MType == gauge:
		ms.InsertMetric(metric.ID, *metric.Value)
	case metric.MType == counter:
		ms.CountCounterMetric(metric.ID, *metric.Delta)
	default:
		err := fmt.Errorf("this type of metric doesnt't exist")
		return err
	}
	return nil
}

func (ms *MetricsStorage) GetAllMetrics() []encoder.Encode {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var allMetrics []encoder.Encode

	for k, v := range ms.MetricsGauge {
		var metric encoder.Encode
		v = ms.MetricsGauge[k]
		metric.MType = "gauge"
		metric.ID = k
		metric.Value = &v
		allMetrics = append(allMetrics, metric)
	}

	for k, d := range ms.MetricsCounter {
		var metric encoder.Encode
		d = ms.MetricsCounter[k]
		metric.MType = "gauge"
		metric.ID = k
		metric.Delta = &d
		allMetrics = append(allMetrics, metric)
	}

	return allMetrics
}

func (ms *MetricsStorage) GetBackupCounter(id string, delta int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.MetricsCounter[id] = delta
}

func (ms *MetricsStorage) GetBackupGauge(id string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.MetricsGauge[id] = value
}

func (ms *MetricsStorage) Load() error {
	metrics, err := ms.Keeper.LoadMetrics()
	if err != nil {
		return err
	}

	for _, v := range metrics {
		switch v.MType {
		case "gauge":
			ms.GetBackupGauge(v.ID, *v.Value)
		case "counter":
			ms.GetBackupCounter(v.ID, *v.Delta)
		}
	}
	return nil
}
