package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

type metricStorage struct {
	storage storage.MemeS
}

func NewMetricStorage(s storage.MemeS) metricStorage {
	return metricStorage{storage: s}
}

func (ms metricStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

func (ms metricStorage) AllMetrics(w http.ResponseWriter, r *http.Request) {

	allMetrics := ms.storage.GetAllMetrics()
	Form := fmt.Sprintf("%s", allMetrics)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Form))
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
