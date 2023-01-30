package encoder

import (
	"encoding/json"
	"fmt"
	"io"
)

type Encode struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func New() *Encode {
	return &Encode{}
}

func (m *Encode) Decode(body io.Reader) error {
	dec := json.NewDecoder(body)

	if err := dec.Decode(&m); err != nil {
		return err
	}

	return nil
}

func (m *Encode) Encode(body io.Writer) error {
	enc := json.NewEncoder(body)

	if err := enc.Encode(&m); err != nil {
		return err
	}

	return nil
}

func (m *Encode) Marshall() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Encode) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return fmt.Errorf("error during unmarshalling: %w", err)
	}

	return nil
}
