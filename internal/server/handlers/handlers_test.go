package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/popooq/collectimg-ma/internal/utils/coder"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

func NewRouter() chi.Router {

	MemS := storage.NewMemStorage()
	metricStruct := coder.NewMetricsStruct()
	handler := NewMetricStorage(MemS, metricStruct)

	MemS.InsertMetric("Alloc", 123.000)
	MemS.CountCounterMetric("PollCount", 34)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/update/{mType}/{mName}/{mValue}", func(w http.ResponseWriter, r *http.Request) {
			handler.CollectMetrics(w, r)
		})
		r.Get("/value/{mType}/{mName}", func(w http.ResponseWriter, r *http.Request) {
			handler.MetricValue(w, r)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handler.AllMetrics(w, r)
		})
	})
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
			name: "Positive test Counter",
			url:  "/value/counter/PollCount",
			code: 200,
		},
		{
			name: "Negative test Gauge",
			url:  "/value/gauge/Allo",
			code: 404,
		},
		{
			name: "Negative test Counter",
			url:  "/value/counter/PollCoun",
			code: 404,
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
