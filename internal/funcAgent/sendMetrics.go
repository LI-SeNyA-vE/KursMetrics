package funcagent

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/errorRetriable"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

//// SendingMetric Функция которая каджые $FlagPollInterval секунд запускает функию по отправке метрик
//func SendingMetric(gaugeMetrics map[string]float64, counterMetrics map[string]int64, flagPollInterval int64, flagReportInterval int64, flagAddressAndPort string) {
//	ticker1 := time.NewTicker(time.Duration(flagPollInterval) * time.Second)
//	ticker2 := time.NewTicker(time.Duration(flagReportInterval) * time.Second)
//	defer ticker1.Stop()
//	defer ticker2.Stop()
//
//	for {
//		select {
//		case <-ticker1.C:
//			gaugeMetrics, counterMetrics = UpdateMetric()
//			fmt.Printf("Пауза в %d секунд между сборкой метрик\n", flagPollInterval)
//		case <-ticker2.C:
//			SendJSONMetricsGauge(gaugeMetrics, flagAddressAndPort)
//			SendJSONMetricsCounter(counterMetrics, flagAddressAndPort)
//			fmt.Printf("Пауза в %d секунд между отправкой метрик на сервер\n", flagReportInterval)
//		}
//	}
//}

// SendJSONMetricsGauge Отправляет метрики типа Gauge по по url
//func SendJSONMetricsGauge(mapMetric map[string]float64, flagAddressAndPort string) {
//
//	client := resty.New()
//	url := fmt.Sprintf("http://%s/update/", flagAddressAndPort)
//
//	for nameMetric, value := range mapMetric {
//		metrics := metricStorage.Metrics{
//			ID:    nameMetric,
//			MType: "gauge",
//			Value: &value,
//		}
//
//		jsonData, err := json.Marshal(metrics)
//		if err != nil {
//			log.Printf("Ошибка маршалинга метрик в JSON: %v", err)
//			continue
//		}
//
//		compressedData, err := gzipCompress(jsonData)
//		if err != nil {
//			log.Printf("Ошибка сжатия метрик: %v", err)
//			continue
//		}
//
//		_, err = client.R().
//			SetHeader("Content-Type", "application/json").
//			SetHeader("Content-Encoding", "gzip").
//			SetHeader("Accept-Encoding", "gzip").
//			SetBody(compressedData).
//			Post(url)
//		if err != nil {
//			log.Printf("Не удалось отправить метрику %s типа %s с ошибкой: %v", metrics.ID, metrics.MType, err)
//		}
//	}
//
//}

// SendJSONMetricsCounter Отправляет метрики типа Gauge по по url
//func SendJSONMetricsCounter(mapMetric map[string]int64, flagAddressAndPort string) {
//	client := resty.New()
//	url := fmt.Sprintf("http://%s/update/", flagAddressAndPort)
//
//	for nameMetric, value := range mapMetric {
//		metrics := metricStorage.Metrics{
//			ID:    nameMetric,
//			MType: "counter",
//			Delta: &value,
//		}
//
//		jsonData, err := json.Marshal(metrics)
//		if err != nil {
//			log.Printf("Ошибка маршалинга метрик в JSON: %v", err)
//			return
//		}
//
//		compressedData, err := gzipCompress(jsonData)
//		if err != nil {
//			log.Printf("Ошибка сжатия метрик: %v", err)
//			continue
//		}
//
//		_, err = client.R().
//			SetHeader("Content-Type", "application/json").
//			SetHeader("Content-Encoding", "gzip").
//			SetHeader("Accept-Encoding", "gzip").
//			SetBody(compressedData).
//			Post(url)
//		if err != nil {
//			log.Printf("Не удалось отправить метрику %s типа %s с ошибкой: %v", metrics.ID, metrics.MType, err)
//		}
//	}
//}

func SendingBatchMetric(gaugeMetrics map[string]float64, counterMetrics map[string]int64, cfg config.Agent) {
	ticker1 := time.NewTicker(time.Duration(cfg.FlagPollInterval) * time.Second)
	ticker2 := time.NewTicker(time.Duration(cfg.FlagReportInterval) * time.Second)
	defer ticker1.Stop()
	defer ticker2.Stop()

	for {
		select {
		case <-ticker1.C:
			gaugeMetrics, counterMetrics = UpdateMetric()
			fmt.Printf("Пауза в %d секунд между сборкой метрик\n", cfg.FlagPollInterval)
		case <-ticker2.C:
			SendgBatchJSONMetricsGauge(gaugeMetrics, cfg.FlagAddressAndPort, cfg.FlagKey)
			SendgBatchJSONMetricsCounter(counterMetrics, cfg.FlagAddressAndPort, cfg.FlagKey)
			fmt.Printf("Пауза в %d секунд между отправкой 'батчей' метрик на сервер\n", cfg.FlagReportInterval)
		}
	}
}

func SendgBatchJSONMetricsGauge(mapMetric map[string]float64, flagAddressAndPort string, fladKey string) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/updates/", flagAddressAndPort)

	var metrics []metricStorage.Metrics
	for nameMetric, value := range mapMetric {
		metrics = append(metrics, metricStorage.Metrics{
			ID:    nameMetric,
			MType: "gauge",
			Value: value,
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

	_, err = errorRetriable.ErrorRetriableHTTP(func() (interface{}, error) {
		return sendMetrics(client, url, compressedData, fladKey)
	})

	if err != nil {
		log.Printf("Не удалось отправить 'батч' метрик типа 'Gauge' с ошибкой: %v", err)
	}
}

// SendgBatchJSONMetricsCounter Отправляет метрики типа Gauge по по url
func SendgBatchJSONMetricsCounter(mapMetric map[string]int64, flagAddressAndPort string, flagKey string) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/updates/", flagAddressAndPort)

	var metrics []metricStorage.Metrics
	for nameMetric, value := range mapMetric {
		metrics = append(metrics, metricStorage.Metrics{
			ID:    nameMetric,
			MType: "counter",
			Delta: value,
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

	_, err = errorRetriable.ErrorRetriableHTTP(func() (interface{}, error) {
		return sendMetrics(client, url, compressedData, flagKey)
	})

	if err != nil {
		log.Printf("Не удалось отправить 'батч' метрик типа 'Counter' с ошибкой: %v", err)
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

func sendMetrics(client *resty.Client, url string, compressedData []byte, fladKey string) (interface{}, error) {
	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(compressedData)

	if fladKey != "" {
		h := hmac.New(sha256.New, []byte(fladKey))
		h.Write(compressedData)
		hash := hex.EncodeToString(h.Sum(nil))

		request.SetHeader("HashSHA256", hash)
	}
	response, err := request.Post(url)

	// Если произошла ошибка или статус-код не 2xx, возвращаем ошибку
	if err != nil || response.StatusCode() >= 400 {
		return nil, errors.New("ошибка при отправке метрик")
	}

	return response, nil
}
