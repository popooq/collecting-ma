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

func (h *Hash) Hasher(m *encoder.Encode) (string, error) {
	var src string

	if h.Key == nil {
		return "", nil
	}
	switch m.MType {
	case "counter":
		src = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case "gauge":
		src = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}

	hmac := hmac.New(sha256.New, h.Key)
	_, err := hmac.Write([]byte(src))
	if err != nil {
		return "", err
	}
	hash := fmt.Sprintf("%x", hmac.Sum(nil))
	return hash, nil
}

func (h *Hash) HashChecker(hash string, m encoder.Encode) error {
	if m.Hash != "" && !hmac.Equal([]byte(m.Hash), []byte(hash)) {
		return fmt.Errorf("not equal m.hash %x and hash %x", []byte(m.Hash), []byte(hash))
	}
	return nil
}
