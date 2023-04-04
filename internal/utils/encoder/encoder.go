package encoder

import (
	"encoding/json"
	"fmt"
	"io"
)

// Encode хранит информацию о метриках
type Encode struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

// New созщадет новый Encode
func New() *Encode {
	return &Encode{}
}

// Decode декодирует body в структуру Encode
func (m *Encode) Decode(body io.Reader) error {
	dec := json.NewDecoder(body)

	if err := dec.Decode(&m); err != nil {
		return err
	}

	return nil
}

// Encode кодирует содержимое Encode в body
func (m *Encode) Encode(body io.Writer) error {
	enc := json.NewEncoder(body)

	if err := enc.Encode(&m); err != nil {
		return err
	}

	return nil
}

// Marshall сериализует Encode в []byte
func (m *Encode) Marshall() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshall десериализирует data в Encode
func (m *Encode) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return fmt.Errorf("error during unmarshalling: %w", err)
	}

	return nil
}
