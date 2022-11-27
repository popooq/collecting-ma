package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/popooq/collectimg-ma/internal/server/storage"
)

func TestMetricStorageServeHTTP(t *testing.T) {

	tests := []struct {
		name string
		url  string
		code int
	}{
		{
			name: "Positive test: Gauge",
			url:  "/update/gauge/Alloc/12",
			code: 200,
		},
		{
			name: "Positive test: Counter",
			url:  "/update/counter/PollCount/111",
			code: 200,
		},
		{
			name: "Negative test: Unkonwn metric",
			url:  "/update/unknown/poop/111",
			code: 501,
		},
		{
			name: "Negative test: Empty value",
			url:  "/update/counter/Name1",
			code: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()

			memStorage := storage.NewMemStorage()
			handler := NewmetricStorage(memStorage)

			h := http.Handler(handler)
			h.ServeHTTP(w, r)
			result := w.Result()
			if result.StatusCode != tt.code {
				t.Errorf("Expected code %d, got %d", tt.code, result.StatusCode)
			}
			defer result.Body.Close()
		})
	}
}
