package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

	switch m.MType {
	case "counter":
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
		log.Printf("data во время хеширования: %s, дельта: %d", data, *m.Delta)
	case "gauge":
		data = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
		log.Printf("data во время хеширования: %s, значение: %f", data, *m.Value)
	}

	if hsh.Key == nil {
		return ""
	}

	h := hmac.New(sha256.New, hsh.Key)
	h.Write([]byte(data))
	hash := hex.EncodeToString(h.Sum(nil))

	return hash
}

func (hsh *Hash) HashChecker(hash string, m encoder.Encode) error {
	if m.Hash != "" && !hmac.Equal([]byte(m.Hash), []byte(hash)) {
		return fmt.Errorf("not equal m.hash %x and hash %x", []byte(m.Hash), []byte(hash))
	}
	return nil
}
