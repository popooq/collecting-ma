package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/popooq/collectimg-ma/internal/server/storage"
	"github.com/popooq/collectimg-ma/internal/server/trimmer"
)

var form = "%s"

type metricStorage struct {
	storage storage.MemeS
}

func NewmetricStorage(s storage.MemeS) metricStorage {
	return metricStorage{storage: s}
}

func (ms metricStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	fields := trimmer.Trimmer(url)

	switch {
	case fields[1] == "gauge":
		value, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.InsertMetric(fields[2], value)
	case fields[1] == "counter":
		value, err := strconv.Atoi(fields[3])
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.CountCounterMetric(fields[2], uint64(value))
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
	Form := fmt.Sprintf(form, allMetrics)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Form))
}

func (ms metricStorage) MetricValue(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	fmt.Print(mType)
	var mValue string

	switch {
	case mType == "gauge" || mType == "counter":
		value, err := ms.storage.GetMetric(mName)
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
