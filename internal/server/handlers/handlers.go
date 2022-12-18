package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/popooq/collectimg-ma/internal/utils/serializator"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

type (
	metricStorage struct {
		storage      storage.MemeS
		serializator *serializator.Metrics
	}
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "INFO: ", log.Lshortfile)

	infof = func(info string) {
		logger.Output(2, info)
	}
)

func NewMetricStorage(stor storage.MemeS, ser *serializator.Metrics) metricStorage {
	return metricStorage{
		storage:      stor,
		serializator: ser}
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

	err := ms.serializator.Decode(r.Body)
	if err != nil {
		log.Println("something goes wrong")
	}

	switch {
	case ms.serializator.MType == "gauge":
		ms.storage.InsertMetric(ms.serializator.ID, *ms.serializator.Value)
	case ms.serializator.MType == "counter":
		ms.storage.CountCounterMetric(ms.serializator.ID, uint64(*ms.serializator.Delta))
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}
	err = ms.serializator.Encode(w)
	if err != nil {
		log.Println("something goes wrong", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (ms metricStorage) MetricJSONValue(w http.ResponseWriter, r *http.Request) {

	err := ms.serializator.Decode(r.Body)
	if err != nil {
		log.Println("something goes wrong")
	}

	switch {
	case ms.serializator.MType == "gauge":
		gaugeValue, err := ms.storage.GetMetricGauge(ms.serializator.ID)
		if err != nil {
			infof("This metric doesnt exist")
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		ms.serializator.Value = &gaugeValue
	case ms.serializator.MType == "counter":
		value, err := ms.storage.GetMetricCounter(ms.serializator.ID)
		if err != nil {
			infof("This metric doesnt exist")
			http.Error(w, "This metric doesn't exist", http.StatusNotFound)
			return
		}
		counterVal := int64(value)
		ms.serializator.Delta = &counterVal
	default:
		infof("this type of metric doesnt't exist")
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
		return
	}

	err = ms.serializator.Encode(w)
	if err != nil {
		log.Println("simething goes wrong", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
