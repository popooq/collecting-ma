package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/popooq/collectimg-ma/internal/server/handlers"
	"github.com/popooq/collectimg-ma/internal/server/storage"
)

func NewRouter() chi.Router {
	memS := storage.NewMemStorage()
	handler := handlers.NewMetricStorage(memS)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/update/{mType}/{mName}/{mValue}", func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
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
