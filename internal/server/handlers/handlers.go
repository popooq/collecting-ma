package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/popooq/collectimg-ma/internal/utils/coder"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

type (
	metricStorage struct {
		storage storage.MemeS
		coder   *coder.Metrics
	}
)

func NewMetricStorage(stor storage.MemeS, coder *coder.Metrics) metricStorage {
	return metricStorage{
		storage: stor,
		coder:   coder}
}

func (ms metricStorage) CollectMetrics(w http.ResponseWriter, r *http.Request) {

	mTypeParam := chi.URLParam(r, "mType")
	mNameParam := chi.URLParam(r, "mName")
	mValueParam := chi.URLParam(r, "mValue")
	switch {
	case mTypeParam == "gauge":
		value, err := strconv.ParseFloat(mValueParam, 64)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.InsertMetric(mNameParam, value)
	case mTypeParam == "counter":
		value, err := strconv.Atoi(mValueParam)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.CountCounterMetric(mNameParam, uint64(value))
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (ms metricStorage) MetricValue(w http.ResponseWriter, r *http.Request) {

	var mValue string

	mTypeParam := chi.URLParam(r, "mType")
	mNameParam := chi.URLParam(r, "mName")

	switch {
	case mTypeParam == "gauge":
		value, err := ms.storage.GetMetricGauge(mNameParam)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		mValue = fmt.Sprintf("%.3f", value)
	case mTypeParam == "counter":
		value, err := ms.storage.GetMetricCounter(mNameParam)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		mValue = fmt.Sprintf("%d", value)
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mValue))
}

func (ms metricStorage) AllMetrics(w http.ResponseWriter, r *http.Request) {

	allMetrics := ms.storage.GetAllMetrics()
	Form := fmt.Sprintf("%s", allMetrics)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Form))
}

func (ms metricStorage) CollectJSONMetric(w http.ResponseWriter, r *http.Request) {

	err := ms.coder.Decode(r.Body)
	if err != nil {
		log.Println("something goes wrong")
	}

	switch {
	case ms.coder.MType == "gauge":
		ms.storage.InsertMetric(ms.coder.ID, *ms.coder.Value)
		ms.coder.Value, err = ms.storage.GetMetricJSONGauge(ms.coder.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	case ms.coder.MType == "counter":
		ms.storage.CountCounterMetric(ms.coder.ID, uint64(*ms.coder.Delta))
		ms.coder.Delta, err = ms.storage.GetMetricJSONCounter(ms.coder.ID)
		if err != nil {

			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}
	err = ms.coder.Encode(w)
	if err != nil {
		log.Println("something goes wrong", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (ms metricStorage) MetricJSONValue(w http.ResponseWriter, r *http.Request) {

	err := ms.coder.Decode(r.Body)
	if err != nil {
		log.Println("something goes wrong")
	}

	switch {
	case ms.coder.MType == "gauge":
		gaugeValue, err := ms.storage.GetMetricGauge(ms.coder.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		ms.coder.Value = &gaugeValue
	case ms.coder.MType == "counter":
		value, err := ms.storage.GetMetricCounter(ms.coder.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		counterVal := int64(value)
		ms.coder.Delta = &counterVal
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	err = ms.coder.Encode(w)
	if err != nil {
		log.Println("simething goes wrong", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
