package sender

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"strings"

	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

func SendMetrics(value any, name, endpoint, key string) {
	var encoderJSON encoder.Encode
	types := strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%T", value), "storage."))
	encoderJSON.ID = name
	encoderJSON.MType = types
	if encoderJSON.MType == "float64" {
		assertvalue, ok := value.(float64)
		if !ok {
			log.Printf("conversion failed")
		}
		floatvalue := float64(assertvalue)
		encoderJSON.Value = &floatvalue
		encoderJSON.Delta = nil
		encoderJSON.MType = "gauge"
	}
	if encoderJSON.MType == "counter" {
		assertdelta, ok := value.(storage.Counter)
		if !ok {
			log.Printf("conversion failed")
		}
		intdelta := int64(assertdelta)
		encoderJSON.Delta = &intdelta
		encoderJSON.Value = nil
	}
	if key != "" {
		hash, err := hasher.Hasher(encoderJSON, key)
		if err != nil {
			log.Printf("something went wrong %s", err)
		}
		encoderJSON.Hash = hash
	}
	body, err := encoderJSON.Marshall()
	if err != nil {
		log.Printf("error %s in agent", err)
	}
	endpoint, err = url.JoinPath("http://", endpoint, "update/")
	if err != nil {
		log.Printf("url joining failed, error: %s", err)
	}
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Server unreachible, error: %s", err)
	} else {
		defer resp.Body.Close()
	}
}
