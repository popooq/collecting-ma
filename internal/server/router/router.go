package router

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
)

func New(handler handlers.MetricStorage) *chi.Mux {
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
		router.Post("/updates/", func(w http.ResponseWriter, r *http.Request) {
			handler.CollectDBMetrics(w, r)
		})
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handler.AllMetrics(w, r)
		})
		router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			handler.PingDB(w, r)
		})
	})

	return router
}
