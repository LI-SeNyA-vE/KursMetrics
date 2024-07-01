package funcagent

import (
	"fmt"
	"log"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/go-resty/resty/v2"
)

func SendMetricsGauge(mapMetric map[string]Gauge, metricType string) {
	for nameMetric, value := range mapMetric {
		client := resty.New()
		url := fmt.Sprintf("http://%s/update/%s/%s/%f", *config.AddressAndPort, metricType, nameMetric, value)
		_, err := client.R().SetHeader("Content-Type", "text/plain").Post(url)
		if err != nil {
			log.Printf("Не удалось отправить метрику: %s=%f типа %s с ошибкой: %v", nameMetric, value, metricType, err)
		}
	}
}

func SendMetricsCounter(mapMetric map[string]Counter, metricType string) {
	for nameMetric, value := range mapMetric {
		client := resty.New()
		url := fmt.Sprintf("http://%s/update/%s/%s/%d", *config.AddressAndPort, metricType, nameMetric, value)
		_, err := client.R().SetHeader("Content-Type", "text/plain").Post(url)
		if err != nil {
			log.Printf("Не удалось отправить метрику: %s=%d типа %s с ошибкой: %v", nameMetric, value, metricType, err)
		}
	}
}
