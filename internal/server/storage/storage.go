package storage

import "fmt"

type MemeS interface {
	InsertMetric(name string, value any)
	CountCounterMetric(name string, value any)
	GetMetric(name string) (any, error)
	GetAllMetrics() []string
}

type MemStorage struct {
	metrics map[string]any
}

func NewMemStorage() *MemStorage {
	var ms MemStorage
	ms.metrics = make(map[string]any)
	return &ms
}

func (ms *MemStorage) InsertMetric(name string, value any) {
	ms.metrics[name] = value
}

func (ms *MemStorage) CountCounterMetric(name string, value any) {
	val, ok := value.(uint64)
	if ok {
		mapval, ok := ms.metrics[name].(uint64)
		if ok {
			ms.metrics[name] = val + mapval
		}
	}
}

func (ms *MemStorage) GetMetric(name string) (any, error) {
	value, ok := ms.metrics[name]
	if ok {
		return value, nil
	} else {
		err := fmt.Errorf("metric %s doesn't exist", name)
		return nil, err
	}
}

func (ms *MemStorage) GetAllMetrics() []string {
	allMetrics := []string{}
	for k, v := range ms.metrics {
		metric := fmt.Sprintf("%s - %f", k, v)
		allMetrics = append(allMetrics, metric)
	}
	return allMetrics
}
