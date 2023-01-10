package encoder

import (
	"encoding/json"
	"io"
)

type Encode struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func NewEncoderMetricsStruct() *Encode {
	return &Encode{}
}
func (m *Encode) Decode(body io.ReadCloser) error {

	dec := json.NewDecoder(body)

	err := dec.Decode(&m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Encode) Encode(body io.Writer) error {
	enc := json.NewEncoder(body)
	err := enc.Encode(&m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Encode) Marshall() ([]byte, error) {
	return json.Marshal(m)
}
