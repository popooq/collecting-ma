package handlersproto

import (
	"context"
	"fmt"
	"net/http"

	"github.com/popooq/collectimg-ma/internal/storage"
	pb "github.com/popooq/collectimg-ma/proto"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	storage *storage.MetricsStorage
}

func NewMetricServer(storage *storage.MetricsStorage) MetricsServer {
	return MetricsServer{
		storage: storage,
	}
}

func (s *MetricsServer) AddMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var responce pb.AddMetricResponse

	switch {
	case in.Metric.Mtype == gauge:
		s.storage.InsertMetric(in.Metric.ID, in.Metric.Value)
	case in.Metric.Mtype == counter:
		s.storage.CountCounterMetric(in.Metric.ID, in.Metric.Delta)
	}

	return &responce, nil
}

func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var responce pb.GetMetricResponse
	var err error

	switch {
	case in.Mtype == gauge:
		responce.Value, err = s.storage.GetMetricGauge(in.ID)
		if err != nil {
			responce.Error = fmt.Sprintln("This metric doesn't exist", http.StatusNotFound)
		}
	case in.Mtype == counter:
		responce.Delta, err = s.storage.GetMetricCounter(in.ID)
		if err != nil {
			responce.Error = fmt.Sprintln("This metric doesn't exist", http.StatusNotFound)
		}
	}

	return &responce, err
}

func (s *MetricsServer) ListMetric(ctx context.Context, in *pb.ListMetricRequest) (*pb.ListMetricResponse, error) {
	var responce pb.ListMetricResponse
	var err error
	gaugeList, counterList := s.storage.AllMetric()

	for k, v := range gaugeList {
		var metric pb.Metric
		metric.Mtype = gauge
		metric.ID = k
		metric.Value = v
		responce.Metric = append(responce.Metric, &metric)
	}

	for k, d := range counterList {
		var metric pb.Metric
		metric.Mtype = counter
		metric.ID = k
		metric.Delta = d
		responce.Metric = append(responce.Metric, &metric)
	}

	return &responce, err
}
