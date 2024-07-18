package funcagent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"

	config "github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/go-resty/resty/v2"
)

func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("ошибка записи данных в gzip writer: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("ошибка закрытия gzip writer: %w", err)
	}
	return buf.Bytes(), nil
}

func SendJSONMetricsGauge(mapMetric map[string]float64) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/update/", *config.FlagAddressAndPort)

	for nameMetric, value := range mapMetric {
		metrics := metricStorage.Metrics{
			ID:    nameMetric,
			MType: "gauge",
			Value: &value,
		}

		jsonData, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("Ошибка маршалинга метрик в JSON: %v", err)
			continue
		}

		compressedData, err := gzipCompress(jsonData)
		if err != nil {
			log.Printf("Ошибка сжатия метрик: %v", err)
			continue
		}

		_, err = client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Accept-Encoding", "gzip").
			SetBody(compressedData).
			Post(url)
		if err != nil {
			log.Printf("Не удалось отправить метрику %s типа %s с ошибкой: %v", metrics.ID, metrics.MType, err)
		}
	}

}

func SendJSONMetricsCounter(mapMetric map[string]int64) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/update/", *config.FlagAddressAndPort)

	for nameMetric, value := range mapMetric {
		metrics := metricStorage.Metrics{
			ID:    nameMetric,
			MType: "counter",
			Delta: &value,
		}

		jsonData, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("Ошибка маршалинга метрик в JSON: %v", err)
			return
		}

		compressedData, err := gzipCompress(jsonData)
		if err != nil {
			log.Printf("Ошибка сжатия метрик: %v", err)
			continue
		}

		_, err = client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Accept-Encoding", "gzip").
			SetBody(compressedData).
			Post(url)
		if err != nil {
			log.Printf("Не удалось отправить метрику %s типа %s с ошибкой: %v", metrics.ID, metrics.MType, err)
		}
	}
}
