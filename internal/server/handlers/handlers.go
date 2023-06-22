// пакет handlers содержит обработчики сервера
//
// хочу тут еще пару строк
//
// типа важный пакет
//
// типа много инфы
package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/encryptor"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

type (
	// Handler содержит информацию о хендлере
	Handler struct {
		storage       *storage.MetricsStorage
		hasher        *hasher.Hash
		encryptor     string
		trustedSubnet string
	}
	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
)

// New создает новый хендлер
func New(stor *storage.MetricsStorage, hasher *hasher.Hash, restore bool, tsubnet, enc string) Handler {
	if restore {
		err := stor.Load()
		if err != nil {
			log.Printf("error during load from file %s", err)
		}
	}
	return Handler{
		storage:       stor,
		hasher:        hasher,
		encryptor:     enc,
		trustedSubnet: tsubnet,
	}

}

// Route создает новый роутер
func (h Handler) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(h.trustedWare)

	r.Post("/update/{mType}/{mName}/{mValue}", h.addMetrics)
	r.Get("/value/{mType}/{mName}", h.getMetric)
	r.Post("/update/", h.addJSONMetric)
	r.Post("/value/", h.getMetricJSON)
	r.Post("/updates/", h.addDBMetrics)
	r.Get("/", h.getAllMetrics)
	r.Get("/ping", h.pingDB)

	return r
}

func (h Handler) trustedWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipStr := r.Header.Get("X-Real-IP")
		log.Println(h.trustedSubnet != "" && ipStr != h.trustedSubnet)
		if h.trustedSubnet != "" && ipStr != h.trustedSubnet {
			w.WriteHeader(403)
		}
		next.ServeHTTP(w, r)
	})
}

func (h Handler) addMetrics(w http.ResponseWriter, r *http.Request) {
	metricTypeParam := chi.URLParam(r, "mType")
	metricNameParam := chi.URLParam(r, "mName")
	metricValueParam := chi.URLParam(r, "mValue")
	log.Printf("metricTypeParam: %s \n metricNameParam: %s \n metricValueParam: %s", metricTypeParam, metricNameParam, metricValueParam)
	switch {
	case metricTypeParam == gauge:
		value, err := strconv.ParseFloat(metricValueParam, 64)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}

		h.storage.InsertMetric(metricNameParam, value)

	case metricTypeParam == counter:
		value, err := strconv.Atoi(metricValueParam)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}

		h.storage.CountCounterMetric(metricNameParam, int64(value))

	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(nil); err != nil {
		return
	}
}

func (h Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	var metricValue string

	metricTypeParam := chi.URLParam(r, "mType")
	metricNameParam := chi.URLParam(r, "mName")

	switch {
	case metricTypeParam == gauge:
		value, err := h.storage.GetMetricGauge(metricNameParam)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}

		metricValue = fmt.Sprintf("%.3f", value)

	case metricTypeParam == counter:
		value, err := h.storage.GetMetricCounter(metricNameParam)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}

		metricValue = fmt.Sprintf("%d", value)

	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(metricValue))
	if err != nil {
		return
	}
}

func (h Handler) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	allMetrics := h.storage.GetAllMetrics()
	listOfMetrics := fmt.Sprintf("%+v", allMetrics)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(listOfMetrics))
	if err != nil {
		return
	}
}

func (h Handler) addJSONMetric(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error during ReadAll: %s", err)
	}

	encryptor, _ := encryptor.New(h.encryptor, "private")
	body, err = encryptor.Decrypt(body)
	if err != nil {
		log.Panicln(err)
	}
	encoder := encoder.New()

	err = encoder.Unmarshal(body)
	if err != nil {
		log.Printf("error during unmarshalling in handler: %s", err)
	}

	switch {
	case encoder.MType == gauge:
		h.storage.InsertMetric(encoder.ID, *encoder.Value)
		encoder.Value, err = h.storage.GetMetricJSONGauge(encoder.ID)
		encoder.Delta = nil
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	case encoder.MType == counter:
		h.storage.CountCounterMetric(encoder.ID, *encoder.Delta)
		encoder.Delta, err = h.storage.GetMetricJSONCounter(encoder.ID)
		encoder.Value = nil
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	encoder.Hash = h.hasher.Hasher(encoder)

	h.storage.Keeper.SaveMetric(encoder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = encoder.Encode(w)
	if err != nil {
		log.Println("something goes wrong", err)
	}
}

func (h Handler) getMetricJSON(w http.ResponseWriter, r *http.Request) {
	encoder := encoder.New()

	body, err := io.ReadAll(r.Body)
	if err != nil {

		log.Println("read request body error!")
	}

	encryptor, _ := encryptor.New(h.encryptor, "private")
	body, err = encryptor.Decrypt(body)
	if err != nil {
		log.Println(err)
	}

	err = encoder.Unmarshal(body)
	if err != nil {
		http.Error(w, fmt.Sprintln("something went wrong while decoding", err), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	switch {
	case encoder.MType == gauge:
		gaugeValue, err := h.storage.GetMetricGauge(encoder.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)

			return
		}

		encoder.Value = &gaugeValue
		encoder.Delta = nil
	case encoder.MType == counter:
		value, err := h.storage.GetMetricCounter(encoder.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)

			return
		}

		counterVal := value

		encoder.Delta = &counterVal
		encoder.Value = nil

	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	encoder.Hash = h.hasher.Hasher(encoder)

	err = h.hasher.HashChecker(encoder.Hash, *encoder)
	if err != nil {
		http.Error(w, fmt.Sprintf("error : %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = encoder.Encode(w)
	if err != nil {
		log.Println("something went wrong while encode", err)
	}
}

func (h Handler) pingDB(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.Keeper.KeeperCheck(); err != nil {
		http.Error(w, "DataBase doesn't responce", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (h Handler) addDBMetrics(w http.ResponseWriter, r *http.Request) {
	var Metrics []encoder.Encode

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error during ReadAll: %s", err)
	}

	err = json.Unmarshal(body, &Metrics)
	if err != nil {
		log.Printf("error during unmarshalling in handler CollectDBMetrics: %s", err)
	}

	for _, metric := range Metrics {
		err = h.storage.InsertMetrics(metric)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}

		err = h.storage.Keeper.SaveAllMetrics(metric)
		if err != nil {
			log.Printf("error while adding metrics to buffer %s", err)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gzip, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gzip.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzip}, r)
	})
}
