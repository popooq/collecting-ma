package storage

import (
	"fmt"
	"sync"
)

type (
	MemeS interface {
		InsertMetric(name string, value float64)
		CountCounterMetric(name string, value uint64)
		GetMetricGauge(name string) (float64, error)
		GetAllMetrics() []string
		GetMetricCounter(name string) (uint64, error)
	}

	memStorage struct {
		metricsGauge   map[string]float64
		metricsCounter map[string]uint64
		mu             sync.Mutex
	}
)

func NewMemStorage() *memStorage {
	var ms memStorage
	ms.mu = sync.Mutex{}
	ms.metricsGauge = make(map[string]float64)
	ms.metricsCounter = make(map[string]uint64)
	return &ms
}

func (ms *memStorage) InsertMetric(name string, value float64) {
	ms.mu.Lock()
	ms.metricsGauge[name] = value
	ms.mu.Unlock()
}

func (ms *memStorage) CountCounterMetric(name string, value uint64) {
	ms.mu.Lock()
	ms.metricsCounter[name] += value
	ms.mu.Unlock()
}

func (ms *memStorage) GetMetricGauge(name string) (float64, error) {
	ms.mu.Lock()
	value, ok := ms.metricsGauge[name]
	ms.mu.Unlock()
	if ok {
		return value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return 0, err
	}
}

func (ms *memStorage) GetMetricCounter(name string) (uint64, error) {
	ms.mu.Lock()
	value, ok := ms.metricsCounter[name]
	ms.mu.Unlock()
	if ok {
		return value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return 0, err
	}
}
func (ms *memStorage) GetAllMetrics() []string {
	ms.mu.Lock()
	allMetrics := []string{}
	ms.mu.Unlock()
	for k, v := range ms.metricsGauge {
		metric := fmt.Sprintf("%s - %.3f", k, v)
		allMetrics = append(allMetrics, metric)
	}
	for k, v := range ms.metricsCounter {
		metric := fmt.Sprintf("%s - %d", k, v)
		allMetrics = append(allMetrics, metric)
	}
	return allMetrics
}
