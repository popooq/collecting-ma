package sender

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/popooq/collectimg-ma/internal/agent/utils/collector"
	"github.com/popooq/collectimg-ma/internal/utils/coder"
)

func SendMetrics(data collector.MetricsMap) {
	for k, v := range data {
		var coder coder.Metrics
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		types := strings.TrimPrefix(fmt.Sprintf("%T", v), "collector.")
		types = strings.ToLower(types)

		coder.ID = k
		coder.MType = types
		if coder.MType == "gauge" {
			assertv := v.(collector.Gauge)
			floatv := float64(assertv)
			coder.Value = &floatv
		}
		if coder.MType == "counter" {
			assertv, ok := v.(collector.Counter)
			if !ok {
				log.Printf("conversion failed")
			}
			intv := int64(assertv)
			coder.Delta = &intv

		}
		body, err := coder.Marshall()
		if err != nil {
			log.Printf("error %s in agent", err)
		}
		//		log.Printf("&buf %p", &buf)
		log.Printf("struct %+v", coder)
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
