// Package rpchandler реализует "ручки" для работы по rpc
//
// Метод для отправки одной метрики
// SendMetric(context.Context, *SendMetricRequest) (*SendMetricResponse, error)
//
// Метод для отправки нескольких метрик (batch)
// SendBatchMetrics(context.Context, *BatchMetricsRequest) (*BatchMetricsResponse, error)
//
// Метод для получения метрики по имени
// GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
//
// Метод для получения всех метрик
// GetAllMetrics(context.Context, *Empty) (*AllMetricsResponse, error)
package rpchandler

import (
	"context"
	"fmt"
	pb "github.com/LI-SeNyA-vE/KursMetrics/api/proto/v1/metrics"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/storages"
	"github.com/sirupsen/logrus"
	_ "google.golang.org/grpc"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
	storage storages.MetricsStorage
	log     *logrus.Entry
}

func NewMetricsServer(storage storages.MetricsStorage, log *logrus.Entry) *MetricsServer {
	return &MetricsServer{
		storage: storage,
		log:     log,
	}
}

func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var response pb.GetMetricResponse
	response.Metric.Id = in.Id
	response.Metric.Type = in.Type

	switch in.Type {
	case pb.MetricType_GAUGE:
		gauge, err := s.storage.GetGauge(in.Id)
		if err != nil {
			s.log.Errorf("при запросе к БД на GetGauge, ошибка %v", err)
			return nil, err
		}
		response.Metric.Value = gauge

	case pb.MetricType_COUNTER:
		counter, err := s.storage.GetCounter(in.Id)
		if err != nil {
			s.log.Errorf("при запросе к БД на GetCounter, ошибка %v", err)
			return nil, err
		}
		response.Metric.Delta = counter

	case pb.MetricType_UNKNOWN:
		err := fmt.Errorf("в запросе передан неверный тип метрики")
		s.log.Info(err)
		return nil, err

	default:
		err := fmt.Errorf("в запросе передан неверный тип метрики")
		s.log.Info(err)
		return nil, err
	}

	return &response, nil
}

func (s *MetricsServer) GetAllMetrics(ctx context.Context, in *pb.Empty) (*pb.AllMetricsResponse, error) {
	var response pb.AllMetricsResponse
	var allCounter map[string]int64
	var allGauge map[string]float64

	allCounter = s.storage.GetAllCounters()
	allGauge = s.storage.GetAllGauges()

	for id, value := range allGauge {
		response.Metrics = append(response.Metrics, &pb.Metric{
			Id:    id,
			Type:  pb.MetricType_GAUGE,
			Delta: nil,
			Value: &value,
		})
	}

	for id, delta := range allCounter {
		response.Metrics = append(response.Metrics, &pb.Metric{
			Id:    id,
			Type:  pb.MetricType_COUNTER,
			Delta: &delta,
			Value: nil,
		})
	}

	return &response, nil
}
