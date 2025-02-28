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
	"errors"
	"fmt"
	pb "github.com/LI-SeNyA-vE/KursMetrics/api/proto/v1/metrics"
	"strings"
)

func (s *MetricsServer) SendMetric(ctx context.Context, in *pb.SendMetricRequest) (*pb.SendMetricResponse, error) {
	var response pb.SendMetricResponse

	switch in.Metric.Type {
	case pb.MetricType_GAUGE:
		s.storage.UpdateGauge(in.Metric.Id, *in.Metric.Value)

	case pb.MetricType_COUNTER:
		s.storage.UpdateCounter(in.Metric.Id, *in.Metric.Delta)

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

func (s *MetricsServer) SendBatchMetrics(ctx context.Context, in *pb.BatchMetricsRequest) (*pb.BatchMetricsResponse, error) {
	var response pb.BatchMetricsResponse
	var errs strings.Builder
	errs.WriteString("не обработан запрос: ")

	for _, metric := range in.Metrics {
		switch metric.Type {
		case pb.MetricType_GAUGE:
			s.storage.UpdateGauge(metric.Id, *metric.Value)

		case pb.MetricType_COUNTER:
			s.storage.UpdateCounter(metric.Id, *metric.Delta)

		case pb.MetricType_UNKNOWN:
			errs.WriteString(fmt.Sprint("метрика - ", metric.Id, " с типом: ", metric.Type, ", "))

		default:
			errs.WriteString(fmt.Sprint("метрика - ", metric.Id, " с типом: ", metric.Type, ", "))

		}
	}

	if errs.Len() > len("не обработан запрос: ") {
		s.log.Error(errs.String())
		return &response, errors.New(errs.String())
	}

	return &response, nil
}
