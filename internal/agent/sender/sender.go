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

	//log.Printf("metric: %s, it's hash: %s", encoderJSON.ID, encoderJSON.Hash)

	body, err := encoderJSON.Marshall()
	if err != nil {
		log.Printf("error %s in agent", err)
	}

	requestBody := bytes.NewBuffer(body)

	endpoint, err = url.JoinPath("http://", endpoint, "update/")
	if err != nil {
		log.Printf("url joining failed, error: %s", err)
	}

	log.Println(string(body))

	// switch encoderJSON.MType {
	// case "gauge":
	// 	log.Printf("json в разрезре: Имя %s \n Тип %s \n Значение %f \n Хеш %s", encoderJSON.ID, encoderJSON.MType, *encoderJSON.Value, encoderJSON.Hash)
	// case "counter":
	// 	log.Printf("json в разрезре: Имя %s \n Тип %s \n Дельта %d \n Хеш %s", encoderJSON.ID, encoderJSON.MType, *encoderJSON.Delta, encoderJSON.Hash)
	// }
	resp, err := http.Post(endpoint, "application/json", requestBody)
	if err != nil {
		log.Printf("Server unreachible, error: %s", err)
	} else {
		defer resp.Body.Close()
	}
}
