package main

import (
	"fmt"

	"time"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	funcAgent "github.com/LI-SeNyA-vE/KursMetrics/internal/funcAgent"
)

func main() {
	cfg := config.GetConfig()
	config.InitializeGlobals(cfg)

	var metrics funcAgent.SystemMetrics
	gaugeMetrics, counterMetrics := metrics.UpdateMetric()

	ticker1 := time.NewTicker(time.Duration(*config.PollInterval) * time.Second)
	ticker2 := time.NewTicker(time.Duration(*config.RreportInterval) * time.Second)
	for {
		select {
		case <-ticker1.C:
			metrics.UpdateMetric()
			fmt.Printf("Пауза в %d секунд между сборкой метрик\n", *config.PollInterval)
		case <-ticker2.C:
			funcAgent.SendMetricsGauge(gaugeMetrics, "gauge")
			funcAgent.SendMetricsCounter(counterMetrics, "counter")
			fmt.Printf("Пауза в %d секунд между отправкой метрик на сервер\n", *config.RreportInterval)
		}

	}
}