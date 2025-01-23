package send

import (
	"encoding/json"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/utils/errorretriable"
	"github.com/go-resty/resty/v2"
	"log"
)

// SendBatchJSONMetricsCounter Отправляет метрики типа Gauge по по url
func SendBatchJSONMetricsCounter(mapMetric map[string]int64, flagAddressAndPort string, flagKey string) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/updates/", flagAddressAndPort)

	var metrics []storages.Metrics
	for nameMetric, value := range mapMetric {
		metrics = append(metrics, storages.Metrics{
			ID:    nameMetric,
			MType: "counter",
			Delta: &value,
		})
	}

	if metrics == nil {
		return
	}

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("Ошибка маршалинга метрик в JSON: %v", err)
	}

	compressedData, err := gzipCompress(jsonData)
	if err != nil {
		log.Printf("Ошибка сжатия метрик: %v", err)
	}

	_, err = errorretriable.ErrorRetriableHTTP(func() (interface{}, error) {
		return sendMetrics(client, url, compressedData, flagKey)
	})

	if err != nil {
		log.Printf("Не удалось отправить 'батч' метрик типа 'Counter' с ошибкой: %v", err)
	}
}
