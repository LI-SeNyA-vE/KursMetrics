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
