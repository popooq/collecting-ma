package sender

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/popooq/collectimg-ma/internal/agent/utils/collector"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

func SendMetrics(metricData collector.MetricsMap) {
	var encoderJSON encoder.Metrics
	for k, v := range metricData {
		types := strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%T", v), "collector."))
		log.Printf("JSON before : %+v", encoderJSON)
		encoderJSON.ID = k
		encoderJSON.MType = types
		if encoderJSON.MType == "gauge" {
			assertvalue, ok := v.(collector.Gauge)
			if !ok {
				log.Printf("conversion failed")
			}
			floatvalue := float64(assertvalue)
			encoderJSON.Value = &floatvalue
			encoderJSON.Delta = nil
		}
		if encoderJSON.MType == "counter" {
			assertdelta, ok := v.(collector.Counter)
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
		log.Printf("JSOM after: %+v", encoderJSON)
		endpoint := "http://127.0.0.1:8080/update/"
		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Server unreachible")
		}
		defer resp.Body.Close()
	}
}
