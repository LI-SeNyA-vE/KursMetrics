//Package send содержит функции для отправки метрик с Агента на Сервер.
//SendBatchJSONMetrics собирает метрики обоих типов (gauge и counter),
//сериализует их в JSON, сжимает (gzip) и отправляет на соответствующий адрес.
//Также поддерживается повторная отправка (retry) при возникновении ошибок.

package send

import (
	"encoding/json"
	"fmt"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/utils/errorretriable"
	"github.com/go-resty/resty/v2"
	"log"
)

// SendBatchJSONMetrics получает карты gauge- и counter-метрик, конструирует из них срез Metrics,
// преобразует в JSON, сжимает gzip'ом и делает POST-запрос на эндпоинт /updates/.
// Если указано значение flagKey, добавляется HMAC SHA256.
// В случае ошибок отправки использует повторные попытки (errorretriable).
func SendBatchJSONMetrics(mapMetricGauge map[string]float64, mapMetricCounter map[string]int64, flagAddressAndPort string, flagKey string) {
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
		return sendMetrics(client, url, compressedData, flagKey)
	})

	if err != nil {
		log.Printf("Не удалось отправить 'батч' метрик с ошибкой: %v", err)
	} else {
		log.Printf("Отправили метрики. Gauge=%d, Counter=%d\n",
			len(mapMetricGauge),
			len(mapMetricGauge),
		)
	}
}
