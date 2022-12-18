package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/utils/coder"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

func NewRouter() chi.Router {
	memS := storage.NewMemStorage()
	metricStruct := coder.NewMetricsStruct()
	handler := handlers.NewMetricStorage(memS, metricStruct)

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
		r.Post("/update/", func(w http.ResponseWriter, r *http.Request) {
			handler.CollectJSONMetric(w, r)
		})
		r.Post("/value/", func(w http.ResponseWriter, r *http.Request) {
			handler.MetricJSONValue(w, r)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handler.AllMetrics(w, r)
		})
	})
	return r
}
