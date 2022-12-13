package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

type (
	Metrics struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	}
	metricStorage struct {
		storage storage.MemeS
	}
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "INFO: ", log.Lshortfile)

	infof = func(info string) {
		logger.Output(2, info)
	}
)

func NewMetricStorage(s storage.MemeS) metricStorage {
	return metricStorage{storage: s}
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

	var m, nm Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch {
	case m.MType == "gauge":
		newValue, err := ms.storage.GetMetricGauge(m.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		nm.Value = &newValue
		ms.storage.InsertMetric(m.ID, *m.Value)
	case m.MType == "counter":
		newValue, err := ms.storage.GetMetricCounter(m.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		newDelta := int64(newValue)
		nm.Delta = &newDelta
		ms.storage.CountCounterMetric(m.ID, uint64(*m.Delta))
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}
	nm.ID = m.ID
	nm.MType = m.MType
	/**switch {
	case nm.MType == "gauge":
		newValue, err := ms.storage.GetMetricGauge(m.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		nm.Value = &newValue
	case nm.MType == "counter":
		newValue, err := ms.storage.GetMetricCounter(m.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		newDelta := int64(newValue)
		nm.Delta = &newDelta
	default:
		http.Error(w, "this type of metric doesn't exist", http.StatusNotImplemented)
	}
	**/
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nm)
}

func (ms metricStorage) MetricJSONValue(w http.ResponseWriter, r *http.Request) {

	var m, nm Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nm.MType = m.MType
	nm.ID = m.ID
	switch {
	case m.MType == "gauge":
		gaugeValue, err := ms.storage.GetMetricGauge(m.ID)
		if err != nil {
			infof("This metric doesnt exist")
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		nm.Value = &gaugeValue
	case m.MType == "counter":
		value, err := ms.storage.GetMetricCounter(m.ID)
		if err != nil {
			infof("This metric doesnt exist")
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		counterVal := int64(value)
		nm.Delta = &counterVal
	default:
		infof("this type of metric doesnt't exist")
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nm)
}
