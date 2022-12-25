package storage

import (
	"fmt"
	"sync"
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
		GetBackupCounter(id string, delta int64)
		GetBackupGauge(id string, delta float64)
	}

	MetricsStorage struct {
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
		mu             *sync.Mutex
	}
	Gauge   float64
	Counter int64
)

func NewMetricStorage() *MetricsStorage {
	var ms MetricsStorage
	ms.mu = &sync.Mutex{}
	ms.MetricsGauge = make(map[string]float64)
	ms.MetricsCounter = make(map[string]int64)
	return &ms
}

func (ms *MetricsStorage) InsertMetric(name string, value float64) {
	ms.mu.Lock()
	ms.MetricsGauge[name] = value
	ms.mu.Unlock()
}

func (ms *MetricsStorage) CountCounterMetric(name string, value int64) {
	ms.mu.Lock()
	ms.MetricsCounter[name] += value
	ms.mu.Unlock()
}

func (ms *MetricsStorage) GetMetricGauge(name string) (float64, error) {
	ms.mu.Lock()
	value, ok := ms.MetricsGauge[name]
	ms.mu.Unlock()
	if ok {
		return value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return 0, err
	}
}

func (ms *MetricsStorage) GetMetricJSONGauge(name string) (*float64, error) {
	ms.mu.Lock()
	value, ok := ms.MetricsGauge[name]
	ms.mu.Unlock()
	if ok {
		return &value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return nil, err
	}
}

func (ms *MetricsStorage) GetMetricCounter(name string) (int64, error) {
	ms.mu.Lock()
	value, ok := ms.MetricsCounter[name]
	ms.mu.Unlock()
	if ok {
		return value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return 0, err
	}
}

func (ms *MetricsStorage) GetMetricJSONCounter(name string) (*int64, error) {
	ms.mu.Lock()
	uvalue, ok := ms.MetricsCounter[name]
	ms.mu.Unlock()
	if ok {
		value := int64(uvalue)
		return &value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return nil, err
	}
}
func (ms *MetricsStorage) GetAllMetrics() []string {
	ms.mu.Lock()
	allMetrics := []string{}
	ms.mu.Unlock()
	for k, v := range ms.MetricsGauge {
		metric := fmt.Sprintf("%s - %.3f", k, v)
		allMetrics = append(allMetrics, metric)
	}
	for k, v := range ms.MetricsCounter {
		metric := fmt.Sprintf("%s - %d", k, v)
		allMetrics = append(allMetrics, metric)
	}
	return allMetrics
}

func (ms *MetricsStorage) GetBackupCounter(id string, delta int64) {
	ms.mu.Lock()
	ms.MetricsCounter[id] = delta
	ms.mu.Unlock()
}

func (ms *MetricsStorage) GetBackupGauge(id string, value float64) {
	ms.mu.Lock()
	ms.MetricsGauge[id] = value
	ms.mu.Unlock()
}
