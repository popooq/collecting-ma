package encoder

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (m *Encode) Hasher(key string) (string, error) {
	var src string
	switch m.MType {
	case "counter":
		src = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
		log.Printf("src: %s", src)
	case "gauge":
		src = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
		log.Printf("src: %s", src)
	}

	bkey := []byte(key)
	h := hmac.New(sha256.New, bkey)
	h.Write([]byte(src))
	hash := fmt.Sprintf("%x", h.Sum(nil))

	if m.Hash != "" && !hmac.Equal([]byte(m.Hash), []byte(hash)) {
		return "", fmt.Errorf("not equal m.hash %x and hash %x", []byte(m.Hash), []byte(hash))
	}

	return hash, nil
}

func (m *Encode) HashChecker(hash string) error {
	if m.Hash != "" && !hmac.Equal([]byte(m.Hash), []byte(hash)) {
		return fmt.Errorf("not equal m.hash %x and hash %x", []byte(m.Hash), []byte(hash))
	}
	return nil
}
