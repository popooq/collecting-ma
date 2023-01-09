package handlers

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

type (
	MetricStorage struct {
		storage *storage.MetricsStorage
		encoder *encoder.Encode
		hasher  *hasher.Hash
	}
	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
)

func NewMetricStorage(stor *storage.MetricsStorage, encoder *encoder.Encode, hasher *hasher.Hash) MetricStorage {
	return MetricStorage{
		storage: stor,
		encoder: encoder,
		hasher:  hasher}
}

func (ms MetricStorage) CollectMetrics(w http.ResponseWriter, r *http.Request) {

	metricTypeParam := chi.URLParam(r, "mType")
	metricNameParam := chi.URLParam(r, "mName")
	metricValueParam := chi.URLParam(r, "mValue")
	switch {
	case metricTypeParam == "gauge":
		value, err := strconv.ParseFloat(metricValueParam, 64)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.InsertMetric(metricNameParam, value)
	case metricTypeParam == "counter":
		value, err := strconv.Atoi(metricValueParam)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.CountCounterMetric(metricNameParam, int64(value))
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (ms MetricStorage) MetricValue(w http.ResponseWriter, r *http.Request) {

	var metricValue string

	metricTypeParam := chi.URLParam(r, "mType")
	metricNameParam := chi.URLParam(r, "mName")

	switch {
	case metricTypeParam == "gauge":
		value, err := ms.storage.GetMetricGauge(metricNameParam)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		metricValue = fmt.Sprintf("%.3f", value)
	case metricTypeParam == "counter":
		value, err := ms.storage.GetMetricCounter(metricNameParam)
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
	w.Write([]byte(metricValue))
}

func (ms MetricStorage) AllMetrics(w http.ResponseWriter, r *http.Request) {

	allMetrics := ms.storage.GetAllMetrics()
	listOfMetrics := fmt.Sprintf("%+v", allMetrics)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(listOfMetrics))
}

func (ms MetricStorage) CollectJSONMetric(w http.ResponseWriter, r *http.Request) {

	err := ms.encoder.Decode(r.Body)
	if err != nil {
		log.Println("something went wrong")
	}
	log.Printf("metrics: %+v", ms.encoder)
	switch {
	case ms.encoder.MType == "gauge":
		ms.storage.InsertMetric(ms.encoder.ID, *ms.encoder.Value)

		ms.encoder.Value, err = ms.storage.GetMetricJSONGauge(ms.encoder.ID)
		ms.encoder.Delta = nil
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	case ms.encoder.MType == "counter":
		ms.storage.CountCounterMetric(ms.encoder.ID, *ms.encoder.Delta)
		ms.encoder.Delta, err = ms.storage.GetMetricJSONCounter(ms.encoder.ID)
		ms.encoder.Value = nil
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	ms.encoder.Hash = ms.hasher.Hasher(ms.encoder)
	log.Printf("current hash: %s", ms.encoder.Hash)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = ms.encoder.Encode(w)
	if err != nil {
		log.Println("something goes wrong", err)
	}
}

func (ms MetricStorage) MetricJSONValue(w http.ResponseWriter, r *http.Request) {

	err := ms.encoder.Decode(r.Body)
	if err != nil {
		log.Println("something went wrong")
	}

	w.Header().Set("Content-Type", "application/json")

	switch {
	case ms.encoder.MType == "gauge":
		gaugeValue, err := ms.storage.GetMetricGauge(ms.encoder.ID)
		if err != nil {
			log.Printf("this metric doesn't exist %s", ms.encoder.ID)
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}

		ms.encoder.Value = &gaugeValue
		ms.encoder.Delta = nil

	case ms.encoder.MType == "counter":
		value, err := ms.storage.GetMetricCounter(ms.encoder.ID)
		if err != nil {
			log.Printf("this metric doesn't exist %s", ms.encoder.ID)
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		counterVal := int64(value)
		ms.encoder.Delta = &counterVal
		ms.encoder.Value = nil
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	err = ms.hasher.HashChecker(ms.encoder.Hash, *ms.encoder)
	if err != nil {
		http.Error(w, fmt.Sprintf("error : %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = ms.encoder.Encode(w)
	if err != nil {
		log.Println("something went wrong", err)
	}
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
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
