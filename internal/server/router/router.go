package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

func NewRouter() chi.Router {
	memS := storage.NewMetricStorage()
	metricStruct := encoder.NewEncoderMetricsStruct()
	handler := handlers.NewMetricStorage(memS, metricStruct)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(router chi.Router) {
		router.Post("/update/{mType}/{mName}/{mValue}", func(w http.ResponseWriter, r *http.Request) {
			handler.CollectMetrics(w, r)
		})
		router.Get("/value/{mType}/{mName}", func(w http.ResponseWriter, r *http.Request) {
			handler.MetricValue(w, r)
		})
		router.Post("/update/", func(w http.ResponseWriter, r *http.Request) {
			handler.CollectJSONMetric(w, r)
		})
		router.Post("/value/", func(w http.ResponseWriter, r *http.Request) {
			handler.MetricJSONValue(w, r)
		})
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handler.AllMetrics(w, r)
		})
	})
	return router
}
