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

var encoderJSON encoder.Metrics

func SendMetrics(metricData collector.MetricsMap) {
	for k, v := range metricData {

		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		types := strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%T", v), "collector."))

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
		log.Printf("struct %+v", encoderJSON)
		endpoint := "http://127.0.0.1:8080/update/"
		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Server unreachible")
		}
		log.Printf("response header : %+v", resp.Header)
		log.Printf("Resp body %+v", resp.Body)
		defer resp.Body.Close()
	}
}
