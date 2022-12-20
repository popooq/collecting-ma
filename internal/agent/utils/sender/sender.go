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
	for k, v := range metricData {
		var encoder encoder.Metrics
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		types := strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%T", v), "collector."))

		encoder.ID = k
		encoder.MType = types
		if encoder.MType == "gauge" {
			assertvalue, ok := v.(collector.Gauge)
			if !ok {
				log.Printf("conversion failed")
			}
			floatvalue := float64(assertvalue)
			encoder.Value = &floatvalue
		}
		if encoder.MType == "counter" {
			assertdelta, ok := v.(collector.Counter)
			if !ok {
				log.Printf("conversion failed")
			}
			intdelta := int64(assertdelta)
			encoder.Delta = &intdelta

		}
		body, err := encoder.Marshall()
		if err != nil {
			log.Printf("error %s in agent", err)
		}
		log.Printf("struct %+v", encoder)
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
