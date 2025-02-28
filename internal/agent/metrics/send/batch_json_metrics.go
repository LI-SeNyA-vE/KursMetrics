// Package send содержит функции для отправки метрик с Агента на Сервер.
// SendBatchJSONMetricsHTTP собирает метрики обоих типов (gauge и counter),
// сериализует их в JSON, сжимает (gzip) и отправляет на соответствующий адрес.
// Также поддерживается повторная отправка (retry) при возникновении ошибок.
package send

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/LI-SeNyA-vE/KursMetrics/api/proto/v1/metrics"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/storages"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/utils/errorretriable"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// SendBatchJSONMetricsHTTP получает карты gauge- и counter-метрик, конструирует из них срез Metrics,
// преобразует в JSON, сжимает gzip'ом и делает POST-запрос на эндпоинт /updates/.
// Если указано значение flagKey, добавляется HMAC SHA256.
// В случае ошибок отправки использует повторные попытки (errorretriable).
func SendBatchJSONMetricsHTTP(mapMetricGauge map[string]float64, mapMetricCounter map[string]int64, flagAddressAndPort string, flagHashKey, flagRsaKey string) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/updates/", flagAddressAndPort)

	var metrics []storages.Metrics
	// Формируем срез Metrics из карт gauge и counter
	for nameMetric, value := range mapMetricGauge {
		metrics = append(metrics, storages.Metrics{
			ID:    nameMetric,
			MType: "gauge",
			Value: &value,
		})
	}
	for nameMetric, delta := range mapMetricCounter {
		metrics = append(metrics, storages.Metrics{
			ID:    nameMetric,
			MType: "counter",
			Delta: &delta,
		})
	}

	// Если метрик нет, выходим
	if metrics == nil {
		return
	}

	// Сериализуем в JSON
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("Ошибка маршалинга метрик в JSON: %v", err)
	}

	// Сжимаем gzip'ом
	compressedData, err := gzipCompress(jsonData)
	if err != nil {
		log.Printf("Ошибка сжатия метрик: %v", err)
	}

	// Пытаемся отправить с ретраями
	_, err = errorretriable.ErrorRetriableHTTP(func() (interface{}, error) {
		return nil, sendMetricsHTTP(client, url, compressedData, flagHashKey, flagRsaKey)
	})

	if err != nil {
		log.Printf("Не удалось отправить 'батч' метрик с ошибкой: %v", err)
	} else {
		log.Printf("Отправили метрики. Gauge=%d, Counter=%d\n",
			len(mapMetricGauge),
			len(mapMetricCounter),
		)
	}
}

func SendBatchJSONMetricsRPC(mapMetricGauge map[string]float64, mapMetricCounter map[string]int64, flagAddressAndPort string, flagHashKey, flagRsaKey string) {
	var request pb.BatchMetricsRequest
	for id, delta := range mapMetricCounter {
		request.Metrics = append(request.Metrics, &pb.Metric{
			Id:    id,
			Type:  pb.MetricType_COUNTER,
			Delta: &delta,
			Value: nil,
		})
	}
	for id, value := range mapMetricGauge {
		request.Metrics = append(request.Metrics, &pb.Metric{
			Id:    id,
			Type:  pb.MetricType_GAUGE,
			Delta: nil,
			Value: &value,
		})
	}

	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// получаем переменную интерфейсного типа MetricsServiceClient,
	// через которую будем отправлять сообщения
	c := pb.NewMetricsServiceClient(conn)

	_, err = errorretriable.ErrorRetriableHTTP(func() (interface{}, error) {
		return c.SendBatchMetrics(context.Background(), &request)
	})

	if err != nil {
		log.Printf("Не удалось отправить 'батч' метрик с ошибкой: %v", err)
	} else {
		log.Printf("Отправили метрики. Gauge=%d, Counter=%d\n",
			len(mapMetricGauge),
			len(mapMetricCounter),
		)
	}
}
