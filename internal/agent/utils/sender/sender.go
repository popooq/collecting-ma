package sender

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	"net/http"

	//"strconv"
	"strings"

	"github.com/popooq/collectimg-ma/internal/agent/utils/collector"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func SendMetrics(data collector.MetricsMap) {
	for k, v := range data {

		value := fmt.Sprint(v)
		types := strings.TrimPrefix(fmt.Sprintf("%T", v), "collector.")

		endpoint := "http://127.0.0.1:8080/update/" + types + "/" + k + "/" + value
		resp, err := http.Post(endpoint, "text/plain", nil)
		if err != nil {
			fmt.Println("Server unreachible")
		}
		defer resp.Body.Close()
	}
	/** for k, v := range data {

		m := newMetric(k, v)

		metricJSON, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		body := bytes.NewReader(metricJSON)

		endpoint := "http://127.0.0.1:8080/update"
		resp, err := http.Post(endpoint, "application/json", body)
		if err != nil {
			fmt.Println("Server unreachible")
		}

		defer resp.Body.Close()
	} **/
}

/**func newMetric(k string, v any) *Metrics {
	value := fmt.Sprint(v)
	types := strings.TrimPrefix(fmt.Sprintf("%T", v), "collector.")

	var (
		floatValue   float64
		counterValue int64
		m            Metrics
	)
	m.ID = k
	m.MType = types
	if m.MType == "gauge" {
		floatValue, _ = strconv.ParseFloat(value, 64)
		m.Value = &floatValue
	}
	if m.MType != "gauge" {
		fmt.Println(m.MType)
	}
	if m.MType == "counter" {
		cValue, _ := strconv.Atoi(value)
		counterValue = int64(cValue)
		m.Delta = &counterValue
	}
	return &m
}**/
