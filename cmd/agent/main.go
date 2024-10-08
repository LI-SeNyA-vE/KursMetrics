package main

import (
	"fmt"

	"time"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	funcAgent "github.com/LI-SeNyA-vE/KursMetrics/internal/funcAgent"
)

func main() {
	config.InitializeGlobals()
	gaugeMetrics, counterMetrics := funcAgent.UpdateMetric()

	ticker1 := time.NewTicker(time.Duration(*config.FlagPollInterval) * time.Second)
	ticker2 := time.NewTicker(time.Duration(*config.FlagRreportInterval) * time.Second)
	defer ticker1.Stop()
	defer ticker2.Stop()
	go func() {
		for {
			select {
			case <-ticker1.C:
				gaugeMetrics, counterMetrics = funcAgent.UpdateMetric()
				fmt.Printf("Пауза в %d секунд между сборкой метрик\n", *config.FlagPollInterval)
			case <-ticker2.C:
				funcAgent.SendJSONMetricsGauge(gaugeMetrics)
				funcAgent.SendJSONMetricsCounter(counterMetrics)
				fmt.Printf("Пауза в %d секунд между отправкой метрик на сервер\n", *config.FlagRreportInterval)
			}
		}
	}()
	select {}
}
