package serializator

import (
	"encoding/json"
	"io"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetricsStruct() *Metrics {
	return &Metrics{}
}
func (m *Metrics) Decode(body io.ReadCloser) error {

	dec := json.NewDecoder(body)

	err := dec.Decode(&m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Metrics) Encode(body io.Writer) error {
	enc := json.NewEncoder(body)
	err := enc.Encode(&m)
	if err != nil {
		return err
	}
	return nil
}
