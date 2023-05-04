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
	"github.com/popooq/collectimg-ma/internal/utils/encryptor"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

// Sender описывает sender
type Sender struct {
	hasher    *hasher.Hash
	endpoint  string
	encryptor *encryptor.Encryptor
}

// New создает новый Sender
func New(hasher *hasher.Hash, endpoint string, encryptor *encryptor.Encryptor) Sender {
	return Sender{
		hasher:    hasher,
		endpoint:  endpoint,
		encryptor: encryptor,
	}
}

// Go отправляет метрику на сервер
func (s *Sender) Go(value any, name string) {
	body := s.bodyBuild(value, name)

	body, err := s.encryptor.Encrypt(body)
	if err != nil {
		log.Println(err)
	}
	requestBody := bytes.NewBuffer(body)

	endpoint, err := url.JoinPath("http://", s.endpoint, "update/")
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

func (s *Sender) bodyBuild(value any, name string) []byte {
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

		encoderJSON.Value = &assertvalue
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

	return body
}
