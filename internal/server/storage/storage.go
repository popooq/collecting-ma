package storage

import "fmt"

type MemeS interface {
	InsertMetric(name string, value float64)
	CountCounterMetric(name string, value uint64)
	GetMetricGauge(name string) (any, error)
	GetAllMetrics() []string
	GetMetricCounter(name string) (any, error)
}

type MemStorage struct {
	metricsGauge   map[string]float64
	metricsCounter map[string]uint64
}

func NewMemStorage() *MemStorage {
	var ms MemStorage
	ms.metricsGauge = make(map[string]float64)
	ms.metricsCounter = make(map[string]uint64)
	return &ms
}

func (ms *MemStorage) InsertMetric(name string, value float64) {
	ms.metricsGauge[name] = value
}

func (ms *MemStorage) CountCounterMetric(name string, value uint64) {

	ms.metricsCounter[name] += value

}

func (ms *MemStorage) GetMetricGauge(name string) (any, error) {
	value, ok := ms.metricsGauge[name]
	if ok {
		return value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return nil, err
	}
}

func (ms *MemStorage) GetMetricCounter(name string) (any, error) {
	value, ok := ms.metricsCounter[name]
	if ok {
		return value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return nil, err
	}
}
func (ms *MemStorage) GetAllMetrics() []string {
	allMetrics := []string{}
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
