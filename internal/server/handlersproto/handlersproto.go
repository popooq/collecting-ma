package handlersproto

import (
	"context"
	"log"

	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
	pb "github.com/popooq/collectimg-ma/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	gauge   pb.Mtype = 0
	counter pb.Mtype = 1
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	storage *storage.MetricsStorage
	hasher  *hasher.Hash
}

func NewMetricServer(storage *storage.MetricsStorage, hasher *hasher.Hash, restore bool) MetricsServer {
	if restore {
		err := storage.Load()
		if err != nil {
			log.Printf("error during load from file %s", err)
		}
	}
	return MetricsServer{
		storage: storage,
		hasher:  hasher,
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

	in.Metric.Hash = s.hasher.HashergRPC(in.Metric)

	return &responce, nil
}

func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var responce pb.GetMetricResponse
	var err error

	switch {
	case in.Mtype == gauge:
		responce.Value, err = s.storage.GetMetricGauge(in.ID)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "This metric %s doesn't exist", in.ID)
		}
	case in.Mtype == counter:
		responce.Delta, err = s.storage.GetMetricCounter(in.ID)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "This metric %s doesn't exist", in.ID)
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
