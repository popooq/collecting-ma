package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

type keeperMock struct {
}

func (k keeperMock) SaveMetric(metric *encoder.Encode) error {
	return nil
}
func (k keeperMock) SaveAllMetrics(metric encoder.Encode) error {
	return nil
}
func (k keeperMock) LoadMetrics() ([]encoder.Encode, error) {
	return nil, nil
}

func (k keeperMock) KeeperCheck() error {
	return nil
}

func NewRouter() *chi.Mux {
	var keeper keeperMock
	var cfg config.Config
	MemS := storage.New(keeper)
	hasher := hasher.Mew("")
	handler := New(MemS, hasher, cfg.Restore)

	MemS.InsertMetric("Alloc", 123.000)
	MemS.CountCounterMetric("PollCount", 34)

	r := chi.NewRouter()
	r.Mount("/", handler.Route())

	return r
}

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
			name: "Negative test: Gauge",
			url:  "/update/gauge/Alloc/abc",
			code: 400,
		},
		{
			name: "Negative test: Counter",
			url:  "/update/counter/PollCount/avs",
			code: 400,
		},
		{
			name: "Negative test: Unknown metric",
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

			h := NewRouter()
			h.ServeHTTP(w, r)
			result := w.Result()
			if result.StatusCode != tt.code {
				t.Errorf("Expected code %d, got %d", tt.code, result.StatusCode)
			}
			defer result.Body.Close()
		})
	}
}

func TestMetricStorageAllMetrics(t *testing.T) {
	tests := []struct {
		name string
		url  string
		code int
	}{
		{
			name: "Positive test",
			url:  "/",
			code: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			h := NewRouter()
			h.ServeHTTP(w, r)
			result := w.Result()
			if result.StatusCode != tt.code {
				t.Errorf("Expected code %d, got %d", tt.code, result.StatusCode)
			}
			defer result.Body.Close()
		})
	}
}

func TestMetricStorageMetricValue(t *testing.T) {
	tests := []struct {
		name string
		url  string
		code int
	}{
		{
			name: "Positive test Gauge",
			url:  "/value/gauge/Alloc",
			code: 200,
		},
		{
			name: "Negative test Counter",
			url:  "/value/counter/PollCoun",
			code: 404,
		},
		{
			name: "Negative test Gauge",
			url:  "/value/gauge/Allo",
			code: 404,
		},
		{
			name: "Positive test Counter",
			url:  "/value/counter/PollCount",
			code: 200,
		},
		{
			name: "Negative test Unknown type",
			url:  "/value/gg/Alloc",
			code: 501,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			h := NewRouter()
			h.ServeHTTP(w, r)
			result := w.Result()
			if result.StatusCode != tt.code {
				t.Errorf("Expected code %d, got %d", tt.code, result.StatusCode)
			}
			defer result.Body.Close()
		})
	}
}

func TestHandler_collectJSONMetric(t *testing.T) {
	url := "/update/"
	tests := []struct {
		name string
		body string
		code int
	}{
		// TODO: Add test cases.
		{name: "Positive JSON Counter test",
			body: `{"id": "PollCount", "delta": 2345234211616163, "type": "counter"}`,
			code: 200,
		},
		// {
		// 	name: "",
		// 	body: "",
		// 	code: 200,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := bytes.NewBuffer([]byte(tt.body))
			r := httptest.NewRequest(http.MethodPost, url, requestBody)
			w := httptest.NewRecorder()

			h := NewRouter()
			h.ServeHTTP(w, r)
			result := w.Result()
			if result.StatusCode != tt.code {
				t.Errorf("Expected code %d, got %d", tt.code, result.StatusCode)
			}
			defer result.Body.Close()
		})
	}
}
