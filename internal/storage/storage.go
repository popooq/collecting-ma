package storage

import (
	"fmt"
	"sync"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

type (
	// Storage интерфейс обработки метрик
	Storage interface {
		InsertMetric(name string, value float64)
		CountCounterMetric(name string, value int64)
		GetMetricGauge(name string) (float64, error)
		GetAllMetrics() []string
		GetMetricCounter(name string) (int64, error)
		GetMetricJSONGauge(name string) (*float64, error)
		GetMetricJSONCounter(name string) (*int64, error)
		InsertMetrics(metric encoder.Encode) error
		AllMetric() (map[string]float64, map[string]int64)
	}

	// Keeper интерфейс сохранения метрик
	Keeper interface {
		SaveMetric(metric *encoder.Encode) error
		SaveAllMetrics(metric encoder.Encode) error
		LoadMetrics() ([]encoder.Encode, error)
		KeeperCheck() error
	}

	//MetricsStorage
	MetricsStorage struct {
		Keeper         Keeper             // Keeper реализация интерфейса
		MetricsGauge   map[string]float64 // MetricsGauge мапа хранящяя метрики типа gauge
		MetricsCounter map[string]int64   // MetricsCounter мапа хранящяя метрики типа counter
		mu             sync.Mutex
	}
	// Counter счетчик
	Counter int64
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

// New возвращает новый MetricsStorage
func New(Keeper Keeper) *MetricsStorage {
	return &MetricsStorage{
		mu:             sync.Mutex{},
		MetricsGauge:   make(map[string]float64),
		MetricsCounter: make(map[string]int64),
		Keeper:         Keeper,
	}

}

// InsertMetric добавляет новую метрику типа gauge
func (ms *MetricsStorage) InsertMetric(name string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.MetricsGauge[name] = value
}

// CountCounterMetric добавляет новую метрику типа counter
func (ms *MetricsStorage) CountCounterMetric(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.MetricsCounter[name] += value
}

// GetMetricGauge возвращает метрику типа gauge
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

// GetMetricJSONGauge возвращает ссылку на метрику типа gauge
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

// GetMetricCounter возвращает метрику типа counter
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

// GetMetricJSONCounter возвращает ссылку на метрику типа counter
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

// InsertMetrics добавляет метрику в хранилище
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

// GetAllMetrics возвращает все метрики из хранилища
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
		metric.MType = "counter"
		metric.ID = k
		metric.Delta = &d
		allMetrics = append(allMetrics, metric)
	}

	return allMetrics
}

// GetBackupCounter добавляет counter из бекапа
func (ms *MetricsStorage) GetBackupCounter(id string, delta int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.MetricsCounter[id] = delta
}

// GetBackupGauge добавляет gauge из бекапа
func (ms *MetricsStorage) GetBackupGauge(id string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.MetricsGauge[id] = value
}

// Load загружает метрики из бекапа
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

func (ms *MetricsStorage) AllMetric() (map[string]float64, map[string]int64) {
	return ms.MetricsGauge, ms.MetricsCounter
}
