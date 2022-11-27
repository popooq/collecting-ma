package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/popooq/collectimg-ma/internal/server/storage"
	"github.com/popooq/collectimg-ma/internal/utils/trimmer"
)

type metricStorage struct {
	storage storage.MemeS
}

func NewmetricStorage(s storage.MemeS) metricStorage {
	return metricStorage{storage: s}
}

func (ms metricStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	fmt.Println(url)

	fields := trimmer.Trimmer(url)

	fmt.Println(len(fields))

	if len(fields) != 4 {
		http.Error(w, "Wrong address.", http.StatusNotFound)
		return
	}

	switch {
	case fields[1] == "gauge":
		mValue, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.InsertMetric(fields[2], mValue)
	case fields[1] == "counter":
		mValue, err := strconv.Atoi(fields[3])
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.CountCounterMetric(fields[2], uint64(mValue))
	default:
		http.Error(w, "this type of metric doesnt't exist", http.StatusNotImplemented)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
