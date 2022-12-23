package sender

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

var (
	encoderJSON encoder.Metrics
)

// func SendMetricsMap(metricData collector.MetricsMap) {
// 	for k, v := range metricData {
// 		types := strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%T", v), "collector."))
// 		log.Printf("JSON before : %+v", encoderJSON)
// 		encoderJSON.ID = k
// 		encoderJSON.MType = types
// 		if encoderJSON.MType == "gauge" {
// 			assertvalue, ok := v.(storage.Gauge)
// 			if !ok {
// 				log.Printf("conversion failed")
// 			}
// 			floatvalue := float64(assertvalue)
// 			encoderJSON.Value = &floatvalue
// 			encoderJSON.Delta = nil
// 		}
// 		if encoderJSON.MType == "counter" {
// 			assertdelta, ok := v.(storage.Counter)
// 			if !ok {
// 				log.Printf("conversion failed")
// 			}
// 			intdelta := int64(assertdelta)
// 			encoderJSON.Delta = &intdelta
// 			encoderJSON.Value = nil
// 		}
// 		body, err := encoderJSON.Marshall()
// 		if err != nil {
// 			log.Printf("error %s in agent", err)
// 		}
// 		log.Printf("JSON after: %+v", encoderJSON)
// 		endpoint := "http://127.0.0.1:8080/update/"
// 		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
// 		if err != nil {
// 			fmt.Println("Server unreachible")
// 		}
// 		defer resp.Body.Close()
// 	}
// }

func SendMetrics(value any, name string) {
	log.Printf("JSON before : %+v", encoderJSON)
	log.Printf("metric type = %T", value)
	types := strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%T", value), "storage."))
	encoderJSON.ID = name
	encoderJSON.MType = types
	if encoderJSON.MType == "gauge" {
		assertvalue, ok := value.(storage.Gauge)
		if !ok {
			log.Printf("conversion failed")
		}
		floatvalue := float64(assertvalue)
		encoderJSON.Value = &floatvalue
		encoderJSON.Delta = nil
		log.Printf("value %f", *encoderJSON.Value)
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
	body, err := encoderJSON.Marshall()
	if err != nil {
		log.Printf("error %s in agent", err)
	}
	log.Printf("JSON after: %+v", encoderJSON)
	endpoint := "http://127.0.0.1:8080/update/"
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Server unreachible, error: %s", err)
	} else {
		defer resp.Body.Close()
	}
}
