package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

type (
	metricStorage struct {
		storage storage.Storage
		encoder *encoder.Metrics
	}
)

func NewMetricStorage(stor storage.Storage, encoder *encoder.Metrics) metricStorage {
	return metricStorage{
		storage: stor,
		encoder: encoder}
}

func (ms metricStorage) CollectMetrics(w http.ResponseWriter, r *http.Request) {

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
		ms.storage.CountCounterMetric(metricNameParam, uint64(value))
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (ms metricStorage) MetricValue(w http.ResponseWriter, r *http.Request) {

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

func (ms metricStorage) AllMetrics(w http.ResponseWriter, r *http.Request) {

	allMetrics := ms.storage.GetAllMetrics()
	listOfMetrics := fmt.Sprintf("%s", allMetrics)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(listOfMetrics))
}

func (ms metricStorage) CollectJSONMetric(w http.ResponseWriter, r *http.Request) {

	err := ms.encoder.Decode(r.Body)
	if err != nil {
		log.Println("something goes wrong")
	}
	log.Printf("request :%+v", r)
	log.Printf("metric struct before: %+v", ms.encoder)
	switch {
	case ms.encoder.MType == "gauge":
		ms.storage.InsertMetric(ms.encoder.ID, *ms.encoder.Value)
		log.Printf("value addres before: %p", ms.encoder.Value)
		ms.encoder.Value, err = ms.storage.GetMetricJSONGauge(ms.encoder.ID)
		log.Printf("value addres after: %p", ms.encoder.Value)
		ms.encoder.Delta = nil
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	case ms.encoder.MType == "counter":
		ms.storage.CountCounterMetric(ms.encoder.ID, uint64(*ms.encoder.Delta))
		log.Printf("delta addres: %p", ms.encoder.Delta)
		ms.encoder.Delta, err = ms.storage.GetMetricJSONCounter(ms.encoder.ID)
		log.Printf("delta addres: %p", ms.encoder.Delta)
		ms.encoder.Value = nil
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}
	log.Printf("metric struct after: %+v", ms.encoder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = ms.encoder.Encode(w)
	if err != nil {
		log.Println("something goes wrong", err)
	}
	log.Printf("Response %+v", w)
}

func (ms metricStorage) MetricJSONValue(w http.ResponseWriter, r *http.Request) {

	err := ms.encoder.Decode(r.Body)
	if err != nil {
		log.Println("something goes wrong")
	}

	switch {
	case ms.encoder.MType == "gauge":
		gaugeValue, err := ms.storage.GetMetricGauge(ms.encoder.ID)
		log.Printf("counter value = %f", gaugeValue)
		if err != nil {
			log.Printf("this metric doesn't exist %s", ms.encoder.ID)
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		log.Printf("ms.encoder.Value before = %f", *ms.encoder.Value)
		ms.encoder.Value = &gaugeValue
		log.Printf("ms.encoder.Value after = %f", *ms.encoder.Value)
	case ms.encoder.MType == "counter":
		value, err := ms.storage.GetMetricCounter(ms.encoder.ID)
		log.Printf("counter value = %d", value)
		if err != nil {
			log.Printf("this metric doesn't exist %s", ms.encoder.ID)
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		counterVal := int64(value)
		ms.encoder.Delta = &counterVal
		log.Printf("ms.encoder.Delta after = %d", *ms.encoder.Delta)
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = ms.encoder.Encode(w)
	if err != nil {
		log.Println("simething goes wrong", err)
	}
}
