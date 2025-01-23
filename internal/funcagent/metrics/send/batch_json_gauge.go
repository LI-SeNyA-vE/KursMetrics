package send

import (
	"encoding/json"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/utils/errorretriable"
	"github.com/go-resty/resty/v2"
	"log"
)

func SendBatchJSONMetricsGauge(mapMetric map[string]float64, flagAddressAndPort string, fladKey string) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/updates/", flagAddressAndPort)

	var metrics []storages.Metrics
	for nameMetric, value := range mapMetric {
		metrics = append(metrics, storages.Metrics{
			ID:    nameMetric,
			MType: "gauge",
			Value: &value,
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
		return sendMetrics(client, url, compressedData, fladKey)
	})

	if err != nil {
		log.Printf("Не удалось отправить 'батч' метрик типа 'Gauge' с ошибкой: %v", err)
	}
}
