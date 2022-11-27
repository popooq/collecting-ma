package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/popooq/collectimg-ma/internal/server/storage"
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
	fields := strings.Split(url, "/")
	fmt.Println(fields)
	if len(fields) != 4 {
		http.Error(w, "Wrong address.", http.StatusNotFound)
		return
	}
	switch {
	case fields[2] == "guage":
		mValue, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.InsertMetric(fields[2], mValue)
	case fields[2] == "counter":
		mValue, err := strconv.Atoi(fields[4])
		if err != nil {
			http.Error(w, "There is no value", http.StatusBadRequest)
			return
		}
		ms.storage.CountCounterMetric(fields[2], mValue)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
