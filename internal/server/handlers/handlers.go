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

const (
	gauge   string = "gauge"
	counter string = "counter"
)

type (
	MetricStorage struct {
		storage *storage.MetricsStorage
		//	encoder *encoder.Encode
		hasher *hasher.Hash
	}
	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
)

func New(stor *storage.MetricsStorage, encoder *encoder.Encode, hasher *hasher.Hash) MetricStorage {
	return MetricStorage{
		storage: stor,
		hasher:  hasher,
	}
}

func (ms MetricStorage) CollectMetrics(w http.ResponseWriter, r *http.Request) {
	metricTypeParam := chi.URLParam(r, "mType")
	metricNameParam := chi.URLParam(r, "mName")
	metricValueParam := chi.URLParam(r, "mValue")

	switch {
	case metricTypeParam == gauge:
		value, err := strconv.ParseFloat(metricValueParam, 64)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}

		ms.storage.InsertMetric(metricNameParam, value)

	case metricTypeParam == counter:
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
	if _, err := w.Write(nil); err != nil {
		return
	}
}

func (ms MetricStorage) MetricValue(w http.ResponseWriter, r *http.Request) {
	var metricValue string

	metricTypeParam := chi.URLParam(r, "mType")
	metricNameParam := chi.URLParam(r, "mName")

	switch {
	case metricTypeParam == gauge:
		value, err := ms.storage.GetMetricGauge(metricNameParam)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}

		metricValue = fmt.Sprintf("%.3f", value)

	case metricTypeParam == counter:
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

	_, err := w.Write([]byte(metricValue))
	if err != nil {
		return
	}
}

func (ms MetricStorage) AllMetrics(w http.ResponseWriter, r *http.Request) {
	allMetrics := ms.storage.GetAllMetrics()
	listOfMetrics := fmt.Sprintf("%+v", allMetrics)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(listOfMetrics))
	if err != nil {
		return
	}
}

func (ms MetricStorage) CollectJSONMetric(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error during ReadAll: %s", err)
	}

	log.Print(string(body))

	encoder := encoder.New()

	err = encoder.Unmarshal(body)
	if err != nil {
		log.Printf("error during unmarshalling in handler: %s", err)
	}

	switch {
	case encoder.MType == gauge:
		ms.storage.InsertMetric(encoder.ID, *encoder.Value)
		encoder.Value, err = ms.storage.GetMetricJSONGauge(encoder.ID)
		encoder.Delta = nil
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	case encoder.MType == counter:
		ms.storage.CountCounterMetric(encoder.ID, *encoder.Delta)
		encoder.Delta, err = ms.storage.GetMetricJSONCounter(encoder.ID)
		encoder.Value = nil
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	encoder.Hash = ms.hasher.Hasher(encoder)
	log.Printf("current hash: %s", encoder.Hash)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = encoder.Encode(w)
	if err != nil {
		log.Println("something goes wrong", err)
	}
}

func (ms MetricStorage) MetricJSONValue(w http.ResponseWriter, r *http.Request) {
	encoder := encoder.New()

	err := encoder.Decode(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintln("something went wrong while decoding", err), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	switch {
	case encoder.MType == gauge:
		gaugeValue, err := ms.storage.GetMetricGauge(encoder.ID)
		if err != nil {
			log.Printf("this metric doesn't exist %s", encoder.ID)
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)

			return
		}

		encoder.Value = &gaugeValue
		encoder.Delta = nil

	case encoder.MType == counter:
		value, err := ms.storage.GetMetricCounter(encoder.ID)
		if err != nil {
			log.Printf("this metric doesn't exist %s", encoder.ID)
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

	encoder.Hash = ms.hasher.Hasher(encoder)
	log.Printf("current hash: %s", encoder.Hash)

	err = ms.hasher.HashChecker(encoder.Hash, *encoder)
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
