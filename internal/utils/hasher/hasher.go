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

func Mew(key string) *Hash {
	return &Hash{
		Key: []byte(key),
	}
}

func (hsh *Hash) Hasher(metric *encoder.Encode) string {
	var data string

	if hsh.Key == nil {
		return ""
	}

	switch metric.MType {
	case "counter":
		data = fmt.Sprintf("%s:%s:%d", metric.ID, metric.MType, *metric.Delta)
		log.Printf("data во время хеширования: %s, дельта: %d", data, *metric.Delta)
	case "gauge":
		data = fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)
		log.Printf("data во время хеширования: %s, значение: %f", data, *metric.Value)
	}

	h := hmac.New(sha256.New, hsh.Key)
	h.Write([]byte(data))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (hsh *Hash) HashChecker(hash string, metric encoder.Encode) error {
	if metric.Hash != "" && !hmac.Equal([]byte(metric.Hash), []byte(hash)) {
		return fmt.Errorf("not equal metric.hash %x and hash %x", []byte(metric.Hash), []byte(hash))
	}

	return nil
}
