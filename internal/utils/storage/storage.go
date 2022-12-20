package storage

import (
	"fmt"
	"sync"
)

type (
	Storage interface {
		InsertMetric(name string, value float64)
		CountCounterMetric(name string, value uint64)
		GetMetricGauge(name string) (float64, error)
		GetAllMetrics() []string
		GetMetricCounter(name string) (uint64, error)
		GetMetricJSONGauge(name string) (*float64, error)
		GetMetricJSONCounter(name string) (*int64, error)
	}

	metricsStorage struct {
		metricsGauge   map[string]float64
		metricsCounter map[string]uint64
		mu             sync.Mutex
	}
)

func NewMetricStorage() *metricsStorage {
	var ms metricsStorage
	ms.mu = sync.Mutex{}
	ms.metricsGauge = make(map[string]float64)
	ms.metricsCounter = make(map[string]uint64)
	return &ms
}

func (ms *metricsStorage) InsertMetric(name string, value float64) {
	ms.mu.Lock()
	ms.metricsGauge[name] = value
	ms.mu.Unlock()
}

func (ms *metricsStorage) CountCounterMetric(name string, value uint64) {
	ms.mu.Lock()
	ms.metricsCounter[name] += value
	ms.mu.Unlock()
}

func (ms *metricsStorage) GetMetricGauge(name string) (float64, error) {
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

func (ms *metricsStorage) GetMetricJSONGauge(name string) (*float64, error) {
	ms.mu.Lock()
	value, ok := ms.metricsGauge[name]
	ms.mu.Unlock()
	if ok {
		return &value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return nil, err
	}
}

func (ms *metricsStorage) GetMetricCounter(name string) (uint64, error) {
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

func (ms *metricsStorage) GetMetricJSONCounter(name string) (*int64, error) {
	ms.mu.Lock()
	uvalue, ok := ms.metricsCounter[name]
	ms.mu.Unlock()
	if ok {
		value := int64(uvalue)
		return &value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return nil, err
	}
}
func (ms *metricsStorage) GetAllMetrics() []string {
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
