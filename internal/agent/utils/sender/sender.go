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
		var buf bytes.Buffer
		log.Printf("&buf %p", &buf)
		types := strings.TrimPrefix(fmt.Sprintf("%T", v), "collector.")
		types = strings.ToLower(types)

		coder.ID = k
		coder.MType = types
		if coder.MType == "gauge" {
			assertv := v.(collector.Gauge)
			floatv := float64(assertv)
			coder.Value = &floatv
			log.Printf("coder.Value %p, &floatv %p, floatv %f, v %f", coder.Value, &floatv, floatv, v)
		}
		if coder.MType == "counter" {
			assertv, ok := v.(collector.Counter)
			if !ok {
				log.Printf("conversion failed")
			}
			intv := int64(assertv)
			coder.Delta = &intv
			log.Printf("coder.Delta %p, &intv %p , intv %d, v %d ", coder.Delta, &intv, intv, v)
		}
		err := coder.Encode(&buf)
		if err != nil {
			log.Printf("error %s in agent", err)
		}
		//		log.Printf("&buf %p", &buf)
		log.Printf("struct %+v", coder)
		endpoint := "http://127.0.0.1:8080/update/"
		resp, err := http.Post(endpoint, "text/plain", &buf)
		if err != nil {
			fmt.Println("Server unreachible")
		}
		defer resp.Body.Close()
	}
}
