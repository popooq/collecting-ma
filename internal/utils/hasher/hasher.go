package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

type Hash struct {
	Key []byte
}

func MewHash(key string) *Hash {
	return &Hash{
		Key: []byte(key),
	}
}

func (hsh *Hash) Hasher(m *encoder.Encode) string {
	var data string
	if hsh.Key == nil {
		return ""
	}
	switch m.MType {
	case "counter":
		log.Printf("Im in block counter %s, %s, %d", m.ID, m.MType, *m.Delta)
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
		log.Printf("data: %s", data)
	case "gauge":
		log.Printf("I'm in block gauge %s, %s, %f", m.ID, m.MType, *m.Value)
		data = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
		log.Printf("data: %s", data)
	}
	h := hmac.New(sha256.New, hsh.Key)
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (hsh *Hash) HashChecker(hash string, m encoder.Encode) error {
	if m.Hash != "" && !hmac.Equal([]byte(m.Hash), []byte(hash)) {
		return fmt.Errorf("not equal m.hash %x and hash %x", []byte(m.Hash), []byte(hash))
	}
	return nil
}
