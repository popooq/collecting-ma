package sender

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/popooq/collectimg-ma/internal/agent/utils/collector"
)

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
}
