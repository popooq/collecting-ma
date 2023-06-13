package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	pb "github.com/popooq/collectimg-ma/proto"
)

// Hash хранит ключ шифрвания
type Hash struct {
	Key []byte // Key клчю шифрования
}

// Mew создает новый Hash
func New(key string) *Hash {
	return &Hash{
		Key: []byte(key),
	}
}

// Hasher подписывает стерики
func (hsh *Hash) Hasher(metric *encoder.Encode) string {
	var data string

	if hsh.Key == nil {
		return ""
	}

	switch metric.MType {
	case "counter":
		data = fmt.Sprintf("%s:%s:%d", metric.ID, metric.MType, *metric.Delta)
	case "gauge":
		data = fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)
	}

	h := hmac.New(sha256.New, hsh.Key)
	h.Write([]byte(data))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// HashChecker проверяет целостность хеша
func (hsh *Hash) HashChecker(hash string, metric encoder.Encode) error {
	if metric.Hash != "" && !hmac.Equal([]byte(metric.Hash), []byte(hash)) {
		return fmt.Errorf("not equal metric.hash %x and hash %x", []byte(metric.Hash), []byte(hash))
	}

	return nil
}

func (hsh *Hash) HashergRPC(metric *pb.Metric) string {
	var data string

	if hsh.Key == nil {
		return ""
	}

	switch metric.Mtype {
	case "counter":
		data = fmt.Sprintf("%s:%s:%d", metric.ID, metric.Mtype, metric.Delta)
	case "gauge":
		data = fmt.Sprintf("%s:gauge:%f", metric.ID, metric.Value)
	}

	h := hmac.New(sha256.New, hsh.Key)
	h.Write([]byte(data))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (hsh *Hash) HashCheckergRPC(hash string, metric *pb.Metric) error {
	if metric.Hash != "" && !hmac.Equal([]byte(metric.Hash), []byte(hash)) {
		return fmt.Errorf("not equal metric.hash %x and hash %x", []byte(metric.Hash), []byte(hash))
	}

	return nil
}
