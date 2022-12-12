package handlers

import (
	"encoding/json"
	"fmt"
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

func (ms metricStorage) CollectJSONMetric(w http.ResponseWriter, r *http.Request) {

	var m Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch {
	case m.MType == "gauge":
		ms.storage.InsertMetric(m.ID, *m.Value)
		newValue, err := ms.storage.GetMetricGauge(m.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		m.Value = &newValue
	case m.MType == "counter":
		ms.storage.CountCounterMetric(m.ID, uint64(*m.Delta))
		newValue, err := ms.storage.GetMetricCounter(m.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		newDelta := int64(newValue)
		m.Delta = &newDelta
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}
	//buf := bytes.NewBuffer([]byte{})
	//encoder := json.NewEncoder(buf)
	//encoder.Encode(m)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write(buf.Bytes())
	json.NewEncoder(w).Encode(m)
}

func (ms metricStorage) AllMetrics(w http.ResponseWriter, r *http.Request) {

	allMetrics := ms.storage.GetAllMetrics()
	Form := fmt.Sprintf("%s", allMetrics)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Form))
}

func (ms metricStorage) MetricJSONValue(w http.ResponseWriter, r *http.Request) {

	var m Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch {
	case m.MType == "gauge":
		gaugeValue, err := ms.storage.GetMetricGauge(m.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		m.Value = &gaugeValue
	case m.MType == "counter":
		value, err := ms.storage.GetMetricCounter(m.ID)
		if err != nil {
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		counterVal := int64(value)
		m.Delta = &counterVal
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	//buf := bytes.NewBuffer([]byte{})
	//encoder := json.NewEncoder(buf)
	//encoder.Encode(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write(buf.Bytes())
	json.NewEncoder(w).Encode(m)
}
