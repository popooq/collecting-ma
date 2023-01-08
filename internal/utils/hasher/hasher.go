package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

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

func (hash *Hash) Hasher(m *encoder.Encode) string {
	var data string
	if hash.Key == nil {
		return ""
	}
	switch m.MType {
	case "counter":
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case "gauge":
		data = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}
	h := hmac.New(sha256.New, hash.Key)
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (h *Hash) HashChecker(hash string, m encoder.Encode) error {
	if m.Hash != "" && !hmac.Equal([]byte(m.Hash), []byte(hash)) {
		return fmt.Errorf("not equal m.hash %x and hash %x", []byte(m.Hash), []byte(hash))
	}
	return nil
}