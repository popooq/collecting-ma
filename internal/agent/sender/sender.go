package sender

import (
	"bytes"
	"fmt"
	"log"
	"net/url"

	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

type Sender struct {
	hasher *hasher.Hash
}

func NewSender(hasher *hasher.Hash) Sender {
	return Sender{
		hasher: hasher,
	}
}
func (s *Sender) SendMetrics(value any, name, endpoint, key string) {

	var encoderJSON encoder.Encode

	types := strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%T", value), "storage."))

	encoderJSON.ID = name
	encoderJSON.MType = types

	switch encoderJSON.MType {
	case "float64":
		assertvalue, ok := value.(float64)
		if !ok {
			log.Printf("conversion failed")
		}

		floatvalue := float64(assertvalue)

		encoderJSON.Value = &floatvalue
		encoderJSON.MType = "gauge"
	case "counter":
		assertdelta, ok := value.(storage.Counter)
		if !ok {
			log.Printf("conversion failed")
		}

		intdelta := int64(assertdelta)

		encoderJSON.Delta = &intdelta
	}

	hash := s.hasher.Hasher(&encoderJSON)
	err := s.hasher.HashChecker(hash, encoderJSON)
	if err != nil {
		log.Printf("error: %s", err)
	}

	encoderJSON.Hash = hash

	body, err := encoderJSON.Marshall()
	if err != nil {
		log.Printf("error %s in agent", err)
	}

	requestBody := bytes.NewBuffer(body)

	endpoint, err = url.JoinPath("http://", endpoint, "update/")
	if err != nil {
		log.Printf("url joining failed, error: %s", err)
	}

	client := resty.New().SetBaseURL(endpoint)

	req := client.R().
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")

	resp, err := req.SetBody(requestBody).Post(endpoint)
	if err != nil {
		log.Printf("Server unreachible, error: %s", err)
	} else {
		defer resp.RawBody().Close()
	}
}
