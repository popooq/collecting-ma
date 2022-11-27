package storage

import "fmt"

type MemeS interface {
	InsertMetric(name string, value any)
	CountCounterMetric(name string, value any)
	GetMetric(name string) (any, error)
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
	ms.metrics[name] = value.(uint64) + ms.metrics[name].(uint64)
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
