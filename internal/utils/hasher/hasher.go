package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

func Hasher(m encoder.Encode, key string) (string, error) {
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
