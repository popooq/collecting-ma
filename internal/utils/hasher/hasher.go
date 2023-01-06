package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

func Hasher(m encoder.Encode, key string) (string, error) {
	var src string
	switch m.MType {
	case "gauge":
		src = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	case "counter":
		src = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}

	bkey := []byte(key)
	h := hmac.New(sha256.New, bkey)
	h.Write([]byte(src))
	hash := fmt.Sprintf("%x", h.Sum(nil))

	if m.Hash != "" && !hmac.Equal([]byte(m.Hash), []byte(hash)) {
		return "", fmt.Errorf("not equal hash")
	}

	return hash, nil
}
