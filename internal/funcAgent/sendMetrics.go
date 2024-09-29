package funcagent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

// SendingMetric Функция которая каджые $FlagPollInterval секунд запускает функию по отправке метрик
func SendingMetric(gaugeMetrics map[string]float64, counterMetrics map[string]int64, flagPollInterval int64, flagReportInterval int64, flagAddressAndPort string) {
	ticker1 := time.NewTicker(time.Duration(flagPollInterval) * time.Second)
	ticker2 := time.NewTicker(time.Duration(flagReportInterval) * time.Second)
	defer ticker1.Stop()
	defer ticker2.Stop()

	for {
		select {
		case <-ticker1.C:
			gaugeMetrics, counterMetrics = UpdateMetric()
			fmt.Printf("Пауза в %d секунд между сборкой метрик\n", flagPollInterval)
		case <-ticker2.C:
			SendJSONMetricsGauge(gaugeMetrics, flagAddressAndPort)
			SendJSONMetricsCounter(counterMetrics, flagAddressAndPort)
			fmt.Printf("Пауза в %d секунд между отправкой метрик на сервер\n", flagReportInterval)
		}
	}
}

// SendJSONMetricsGauge Отправляет метрики типа Gauge по по url
func SendJSONMetricsGauge(mapMetric map[string]float64, flagAddressAndPort string) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/update/", flagAddressAndPort)

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

// SendJSONMetricsCounter Отправляет метрики типа Gauge по по url
func SendJSONMetricsCounter(mapMetric map[string]int64, flagAddressAndPort string) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/update/", flagAddressAndPort)

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

// gzipCompress архивирует данные для последующего отправления
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
